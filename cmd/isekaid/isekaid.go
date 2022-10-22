package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"text/template"

	"github.com/pborman/uuid"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

var (
	waifudURL         = flag.String("waifud-url", "http://100.78.40.86:23818", "waifud base URL")
	wgPort            = flag.Int("wireguard-sever-port", 28139, "what port to have the kernel listen on wireguard")
	wgHostAddr        = flag.String("wireguard-host-addr", "fc00::da10/112", "what IP range to have for the metadata service")
	wgGuestAddr       = flag.String("wireguard-guest-addr", "fc00::da1a", "the IP address for the metadata service")
	wgInterfaceName   = flag.String("wireguard-interface-name", "waifud-metadata", "the wireguard interface for the metadata service")
	wgHostPrivateKey  = flag.String("wireguard-host-private-key", "./var/waifud-host.privkey", "wireguard host private key path (b64)")
	wgHostPubkey      = flag.String("wireguard-host-public-key", "./var/waifud-host.pubkey", "wireguard host public key path (b64)")
	wgGuestPrivateKey = flag.String("wireguard-guest-private-key", "./var/waifud-guest.privkey", "wireguard guest private key path (b64)")
	wgGuestPubkey     = flag.String("wireguard-guest-public-key", "./var/waifud-guest.pubkey", "wireguard guest public key path (b64)")
)

const mtu = 1280 // min ipv6 MTU

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	go streamSyslog(ctx, "dnsmasq")

	if err := establishMetadataInterface(); err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := loadUAPI(&buf); err != nil {
		log.Fatal(err)
	}

	tun, tnet, err := netstack.CreateNetTUN([]netip.Addr{netip.MustParseAddr(*wgGuestAddr)}, []netip.Addr{}, mtu)
	if err != nil {
		log.Fatal(err)
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(1, "[wg] "))
	dev.IpcSet(buf.String())
	if err := dev.Up(); err != nil {
		log.Fatal(err)
	}

	defer tun.Close()

	lis, err := tnet.ListenTCP(&net.TCPAddr{Port: 80})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/dumpip2mac", func(w http.ResponseWriter, r *http.Request) {
		ip2macLock.RLock()
		defer ip2macLock.RUnlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ip2mac)
	})
	mux.HandleFunc("/instance/", func(w http.ResponseWriter, r *http.Request) {
		sp := strings.Split(r.URL.Path, "/")
		if len(sp) != 4 {
			http.Error(w, "expected /instance/:id/:key", http.StatusBadRequest)
			return
		}
		id := sp[2]
		key := sp[3]
		remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
		log.Printf("id: %s, key: %s, remoteAddr: %s", id, key, remoteAddr)

		ip2macLock.RLock()
		macAddr, ok := ip2mac[remoteAddr]
		ip2macLock.RUnlock()

		if !ok {
			log.Printf("can't find address in ip2mac mapping: %s", remoteAddr)
		}

		// TODO query waifud for the instance metadata
		_ = macAddr

		uid := uuid.Parse(id)
		if uid == nil {
			http.Error(w, "invalid UUID", http.StatusBadRequest)
			return
		}

		u, err := url.Parse(*waifudURL)
		if err != nil {
			http.Error(w, "[unexpected] invalid waifud url??", http.StatusInternalServerError)
			return
		}
		u.Path = path.Join("/api", "cloudinit", id, key)

		r.URL.Path = ""

		resp, err := http.Get(u.String())
		if err != nil {
			http.Error(w, "[unexpected] waifud can't be reached: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for k, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "isekaid: the waifud metadata service - https://github.com/Xe/waifud")
	})

	srv := &http.Server{
		Handler: mux,
	}

	defer func() {
		<-ctx.Done()
		dev.Down()
		lis.Close()
		srv.Close()
	}()

	fmt.Println("isekaid is go!")
	log.Fatal(srv.Serve(lis))
}

func loadUAPI(w io.Writer) error {
	hostPubkey, err := key2hexFromFile(*wgHostPubkey)
	if err != nil {
		return err
	}

	guestPrivateKey, err := key2hexFromFile(*wgGuestPrivateKey)
	if err != nil {
		return err
	}

	templ := template.Must(template.New("wg-quick").Parse(uAPITemplate))
	return templ.Execute(w, struct {
		Privkey  string
		Pubkey   string
		HostPort int
		HostAddr string
	}{
		Privkey:  guestPrivateKey,
		Pubkey:   hostPubkey,
		HostPort: *wgPort,
		HostAddr: *wgHostAddr,
	})
}

const uAPITemplate = `private_key={{.Privkey}}
listen_port=0
public_key={{.Pubkey}}
endpoint=127.0.0.1:{{.HostPort}}
allowed_ip={{.HostAddr}}
persistent_keepalive_interval=30
`

func run(args ...string) error {
	log.Println("running command:", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func establishMetadataInterface() error {
	run("modprobe", "wireguard")
	run("ip", "link", "del", "dev", *wgInterfaceName)
	err := run("ip", "link", "add", "dev", *wgInterfaceName, "type", "wireguard")
	if err != nil {
		return err
	}

	err = run("ip", "address", "add", *wgHostAddr, "dev", *wgInterfaceName)
	if err != nil {
		return err
	}

	err = run("wg", "set", *wgInterfaceName, "private-key", *wgHostPrivateKey, "listen-port", fmt.Sprint(*wgPort))
	if err != nil {
		return err
	}

	err = run("ip", "link", "set", "up", "dev", *wgInterfaceName)
	if err != nil {
		return err
	}

	guestPubKey, err := os.ReadFile(*wgGuestPubkey)
	if err != nil {
		return err
	}

	err = run("wg", "set", *wgInterfaceName, "peer", strings.TrimSpace(string(guestPubKey)), "persistent-keepalive", "25", "allowed-ips", *wgGuestAddr)
	if err != nil {
		return err
	}

	err = run("ip", "route", "replace", *wgGuestAddr+"/32", "dev", *wgInterfaceName, "table", "main")
	if err != nil {
		return err
	}

	return nil
}
