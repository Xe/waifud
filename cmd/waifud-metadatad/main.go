package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Xe/waifud/key2hex"
	"github.com/facebookgo/flagenv"
	"github.com/go-redis/redis/v8"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	redisURL          = flag.String("redis-url", "redis://chrysalis", "the url to dial out to Redis")
	wgPort            = flag.Int("wireguard-sever-port", 28139, "what port to have the kernel listen on wireguard")
	wgHostAddr        = flag.String("wireguard-host-addr", "169.254.169.253/30", "what IP range to have for the metadata service")
	wgGuestAddr       = flag.String("wireguard-guest-addr", "169.254.169.254", "the IP address for the metadata service")
	wgInterfaceName   = flag.String("wireguard-interface-name", "waifud-metadata", "the wireguard interface for the metadata service")
	wgHostPrivateKey  = flag.String("wireguard-host-private-key", "./var/waifud-host.privkey", "wireguard host private key path (b64)")
	wgHostPubkey      = flag.String("wireguard-host-public-key", "./var/waifud-host.pubkey", "wireguard host public key path (b64)")
	wgGuestPrivateKey = flag.String("wireguard-guest-private-key", "./var/waifud-guest.privkey", "wireguard guest private key path (b64)")
	wgGuestPubkey     = flag.String("wireguard-guest-public-key", "./var/waifud-guest.pubkey", "wireguard guest public key path (b64)")
)

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = opname.With(ctx, "main")

	rOptions, err := redis.ParseURL(*redisURL)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("parsing redis url"))
	}

	rdb := redis.NewClient(rOptions)
	defer rdb.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	ms := metadataServer{cli: rdb}
	go ms.Run(ctx)

	select {
	case <-c:
		cancel()
	case <-ctx.Done():
	}
}

func run(args ...string) error {
	log.Println("running command:", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type metadataServer struct {
	cli *redis.Client
}

func (ms metadataServer) Run(ctx context.Context) {
	err := ms.listen(ctx)
	if err != nil {
		ln.FatalErr(ctx, err)
	}
}

func (ms metadataServer) listen(ctx context.Context) error {
	err := establishMetadataInterface()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = loadUAPI(&buf)
	if err != nil {
		return err
	}

	tun, tnet, err := netstack.CreateNetTUN([]net.IP{net.ParseIP(*wgGuestAddr)}, []net.IP{}, 1420)
	if err != nil {
		return err
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(4, ""))
	dev.IpcSet(buf.String())
	dev.Up()

	lis, err := tnet.ListenTCP(&net.TCPAddr{Port: 80})
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		dev.Down()
		lis.Close()
	}()

	mux := http.NewServeMux()
	mux.Handle("/instance/", ms)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "waifud metadata service - https://github.com/Xe/waifud")
	})

	srv := &http.Server{
		Handler: mux,
	}

	ln.FatalErr(ctx, srv.Serve(lis))

	return nil
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

func loadUAPI(w io.Writer) error {
	hostPubkey, err := key2hex.FromFile(*wgHostPubkey)
	if err != nil {
		return err
	}

	guestPrivateKey, err := key2hex.FromFile(*wgGuestPrivateKey)
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
endpoint=192.168.122.1:{{.HostPort}}
allowed_ip={{.HostAddr}}
persistent_keepalive_interval=5
`

func (ms metadataServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	vmID := filepath.Base(dir)

	data, err := ms.cli.Get(r.Context(), vmID+"/"+file).Result()
	if err != nil {
		ln.Error(r.Context(), err, ln.F{"vmID": vmID, "file": file})
		http.Error(w, "not found", http.StatusNotFound)
	}

	fmt.Fprintln(w, data)
}
