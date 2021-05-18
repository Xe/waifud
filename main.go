package main

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
	"github.com/philandstuff/dhall-golang/v5"
)

//go:embed data/* templates/*
var data embed.FS

var (
	distro      = flag.String("distro", "alpine-edge", "the linux distro to install in the VM")
	name        = flag.String("name", "", "the name of the VM, defaults to a random common blade name")
	zvolPrefix  = flag.String("zvol-prefix", "rpool/mkvm-test/", "the prefix to use for zvol names")
	zvolSize    = flag.Int("zvol-size", 25, "the number of gigabytes for the virtual machine disk")
	memory      = flag.Int("memory", 512, "the number of megabytes of ram for the virtual machine")
	cloudConfig = flag.String("user-data", "./var/xe-base.yaml", "path to a cloud-config userdata file")
	useSATA     = flag.Bool("use-sata", false, "use SATA for the VM's disk interface? (needed if using freebsd-12)")
)

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	cdir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("can't find cache dir: %v", err)
	}
	cdir = filepath.Join(cdir, "within", "mkvm")
	os.MkdirAll(filepath.Join(cdir, "nixos"), 0755)
	os.MkdirAll(filepath.Join(cdir, "qcow2"), 0755)
	os.MkdirAll(filepath.Join(cdir, "seed"), 0755)
	vmID := uuid.New().String()

	if *name == "" {
		commonBladeName, err := getName()
		if err != nil {
			log.Fatal(err)
		}
		name = &commonBladeName
	}

	distros, err := getDistros()
	if err != nil {
		log.Fatalf("can't load internal list of distros: %v", err)
	}

	var resultDistro Distro
	var found bool
	qcowPath := filepath.Join(cdir, "nixos", vmID, "nixos.qcow2")

	if *distro == "nixos" {
		found = true
		resultDistro = Distro{
			Name:        "nixos",
			DownloadURL: "file://" + qcowPath,
			Sha256Sum:   "<computed after build>",
			MinSize:     8,
		}
	}

	for _, d := range distros {
		if d.Name == *distro {
			found = true
			resultDistro = d
			if *zvolSize == 0 {
				zvolSize = &d.MinSize
			}
			if *zvolSize < d.MinSize {
				zvolSize = &d.MinSize
			}
		}
	}
	if !found {
		fmt.Printf("can't find distro %s in my list. Here are distros I know about:\n", *distro)
		for _, d := range distros {
			fmt.Println(d.Name)
		}
		os.Exit(1)
	}

	zvol := filepath.Join(*zvolPrefix, *name)
	if resultDistro.Name != "nixos" {
		qcowPath = filepath.Join(cdir, "qcow2", resultDistro.Sha256Sum)
	}

	macAddress, err := randomMac()
	if err != nil {
		log.Fatalf("can't generate mac address: %v", err)
	}

	l, err := connectToLibvirt()
	if err != nil {
		log.Fatalf("can't connect to libvirt: %v", err)
	}

	log.Println("plan:")
	log.Printf("name: %s", *name)
	log.Printf("zvol: %s (%d GB)", zvol, *zvolSize)
	log.Printf("base image url: %s", resultDistro.DownloadURL)
	log.Printf("mac address: %s", macAddress)
	log.Printf("ram: %d MB", *memory)
	log.Printf("id: %s", vmID)
	log.Printf("cloud config: %s", *cloudConfig)
	if *useSATA {
		log.Println("using SATA for the VM disk interface")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press enter if this looks okay: ")
	reader.ReadString('\n')

	_, err = os.Stat(qcowPath)
	if err != nil {
		log.Printf("downloading distro image %s to %s", resultDistro.DownloadURL, qcowPath)
		fout, err := os.Create(qcowPath)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Get(resultDistro.DownloadURL)
		if err != nil {
			log.Fatalf("can't fetch qcow2 for %s (%s): %v", resultDistro.Name, resultDistro.DownloadURL, err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("%s replied %s", resultDistro.DownloadURL, resp.Status)
		}

		_, err = io.Copy(fout, resp.Body)
		if err != nil {
			log.Fatalf("download of %s failed: %v", resultDistro.DownloadURL, err)
		}

		fout.Close()
		resp.Body.Close()

		fin, err := os.Open(qcowPath)
		if err != nil {
			log.Fatal(err)
		}

		hasher := sha256.New()
		if _, err := io.Copy(hasher, fin); err != nil {
			log.Fatal(err)
		}
		hash := hex.EncodeToString(hasher.Sum(nil))

		if hash != resultDistro.Sha256Sum {
			log.Println("hash mismatch, someone is doing something nasty")
			log.Printf("want: %q", resultDistro.Sha256Sum)
			log.Printf("got:  %q", hash)
			os.Exit(1)
		}

		log.Printf("hash check passed (%s)", resultDistro.Sha256Sum)
	}

	tmpl := template.Must(template.ParseFS(data, "templates/*"))
	var buf = bytes.NewBuffer(nil)
	err = tmpl.ExecuteTemplate(buf, "meta-data", struct {
		Name string
		ID   string
	}{
		Name: *name,
		ID:   vmID,
	})
	if err != nil {
		log.Fatalf("can't generate cloud-config: %v", err)
	}

	dir, err := os.MkdirTemp("", "mkvm")
	if err != nil {
		log.Fatalf("can't make directory: %v", err)
	}

	fout, err := os.Create(filepath.Join(dir, "meta-data"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = fout.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fout.Close()

	if *distro != "nixos" {
		err = run("cp", *cloudConfig, filepath.Join(dir, "user-data"))
		if err != nil {
			log.Fatal(err)
		}
	}

	isoPath := filepath.Join(cdir, "seed", fmt.Sprintf("%s-%s.iso", *name, vmID))

	err = run(
		"genisoimage",
		"-output",
		isoPath,
		"-volid",
		"cidata",
		"-joliet",
		"-rock",
		filepath.Join(dir, "meta-data"),
		filepath.Join(dir, "user-data"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ram := *memory * 1024
	buf.Reset()

	// zfs create -V 20G rpool/safe/vm/sena
	err = run("sudo", "zfs", "create", "-V", fmt.Sprintf("%dG", *zvolSize), zvol)
	if err != nil {
		log.Fatalf("can't create zvol %s: %v", zvol, err)
	}

	err = run("sudo", "qemu-img", "convert", "-O", "raw", qcowPath, filepath.Join("/dev/zvol", zvol))
	if err != nil {
		log.Fatalf("can't import qcow2: %v", err)
	}

	err = tmpl.ExecuteTemplate(buf, "base.xml", struct {
		Name       string
		UUID       string
		Memory     int
		ZVol       string
		Seed       string
		MACAddress string
		SATA       bool
	}{
		Name:       *name,
		UUID:       vmID,
		Memory:     ram,
		ZVol:       zvol,
		Seed:       isoPath,
		MACAddress: macAddress,
		SATA:       *useSATA,
	})
	if err != nil {
		log.Fatalf("can't generate VM template: %v", err)
	}

	domain, err := mkVM(l, buf)
	if err != nil {
		log.Printf("can't create domain for %s: %v", *name, err)
		log.Println("you should run this command:")
		log.Println()
		log.Printf("zfs destroy %s", zvol)
		os.Exit(1)
	}

	log.Printf("created %s", domain.Name)
}

func randomMac() (string, error) {
	buf := make([]byte, 6)
	_, err := crand.Read(buf)
	if err != nil {
		return "", err
	}

	buf[0] = (buf[0] | 2) & 0xfe

	return net.HardwareAddr(buf).String(), nil
}

func getName() (string, error) {
	var names []string
	nameData, err := data.ReadFile("data/names.json")
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(nameData, &names)
	if err != nil {
		return "", err
	}

	return names[rand.Intn(len(names))], nil
}

func run(args ...string) error {
	log.Println("running command:", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func connectToLibvirt() (*libvirt.Libvirt, error) {
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("can't dial libvirt: %w", err)
	}

	l := libvirt.New(c)

	_, err = l.AuthPolkit()
	if err != nil {
		return nil, fmt.Errorf("can't auth with polkit: %w", err)
	}

	if err := l.Connect(); err != nil {
		return nil, fmt.Errorf("can't connect: %w", err)
	}

	return l, nil
}

func mkVM(l *libvirt.Libvirt, buf *bytes.Buffer) (*libvirt.Domain, error) {
	domain, err := l.DomainDefineXML(buf.String())
	if err != nil {
		return nil, err
	}
	err = l.DomainCreate(domain)
	return &domain, err
}

type Distro struct {
	Name        string `dhall:"name" json:"name"`
	DownloadURL string `dhall:"downloadURL" json:"download_url"`
	Sha256Sum   string `dhall:"sha256Sum" json:"sha256_sum"`
	MinSize     int    `dhall:"minSize" json:"min_size"`
}

func getDistros() ([]Distro, error) {
	distroData, err := data.ReadFile("data/distros.dhall")
	if err != nil {
		return nil, err
	}

	var result []Distro
	err = dhall.Unmarshal(distroData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
