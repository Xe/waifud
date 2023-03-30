package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/pojntfx/go-nbd/pkg/client"
)

var (
	file      = flag.String("file", "/dev/nbd0", "Path to device file to create")
	raddr     = flag.String("raddr", "pneuma:10809", "Remote address")
	name      = flag.String("name", "default", "Export name")
	list      = flag.Bool("list", false, "List the exports and exit")
	blockSize = flag.Uint("block-size", 4096, "Block size to use; 0 uses the server's preferred block size")
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connected to", conn.RemoteAddr())

	if *list {
		exports, err := client.List(conn)
		if err != nil {
			log.Fatal(err)
		}

		if err := json.NewEncoder(os.Stdout).Encode(exports); err != nil {
			log.Fatal(err)
		}

		return
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for range sigCh {
			if err := client.Disconnect(f); err != nil {
				log.Fatal(err)
			}

			os.Exit(0)
		}
	}()

	if err := client.Connect(conn, f, &client.Options{
		ExportName: *name,
		BlockSize:  uint32(*blockSize),
	}); err != nil {
		log.Fatal(err)
	}
}
