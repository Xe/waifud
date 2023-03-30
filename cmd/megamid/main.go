package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/pojntfx/go-nbd/pkg/backend"
	"github.com/pojntfx/go-nbd/pkg/server"
	"go.etcd.io/bbolt"
)

var (
	file        = flag.String("file", "disk.img", "Path to file to expose")
	addr        = flag.String("addr", ":10809", "Listen address")
	network     = flag.String("network", "tcp", "Listen network (e.g. `tcp` or `unix`)")
	name        = flag.String("name", "default", "Export name")
	description = flag.String("description", "The default export", "Export description")
	readOnly    = flag.Bool("read-only", false, "Whether the export should be read-only")

	bboltDB     = flag.String("bbolt-db", "data", "path for bbolt database")
	bboltBucket = flag.String("bbolt-bucket", "default", "")
)

const (
	blockSize = 4096

	bboltBackendSize = 1073741824 // one GiB
)

type LoggingBackend struct {
	backend.Backend
}

func (lb LoggingBackend) ReadAt(p []byte, off int64) (n int, err error) {
	nBlocks := len(p) / blockSize
	log.Printf("Reading 0x%x bytes from block 0x%x (%d blocks)", len(p), off/4096, nBlocks)

	n, err = lb.Backend.ReadAt(p, off)

	if err != nil {
		log.Printf("error: %v", err)
	}

	return n, err
}

func (lb LoggingBackend) WriteAt(p []byte, off int64) (n int, err error) {
	nBlocks := len(p) / blockSize
	log.Printf("Writing 0x%x bytes to block 0x%x (%d blocks)", len(p), off/4096, nBlocks)

	n, err = lb.Backend.WriteAt(p, off)

	if err != nil {
		log.Printf("error: %v", err)
	}

	return n, err
}

func main() {
	flag.Parse()

	l, err := net.Listen(*network, *addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.Println("Listening on", l.Addr())

	db, err := bbolt.Open(*bboltDB, 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})

	var b backend.Backend = LoggingBackend{&BoltBackend{
		BucketName: *bboltBucket,
		SizeBytes:  bboltBackendSize,

		db: db,
	}}

	var f *os.File
	if *readOnly {
		f, err = os.OpenFile(*file, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.OpenFile(*file, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()

	b = backend.NewFileBackend(f)

	clients := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Could not accept connection, continuing:", err)

			continue
		}

		clients++

		log.Printf("%v clients connected", clients)

		go func() {
			defer func() {
				_ = conn.Close()

				clients--

				if err := recover(); err != nil {
					log.Printf("Client disconnected with error: %v", err)
				}

				log.Printf("%v clients connected", clients)
			}()

			if err := server.Handle(
				conn,
				[]server.Export{
					{
						Name:        *name,
						Description: *description,
						Backend:     LoggingBackend{b},
					},
				},
				&server.Options{
					ReadOnly:           *readOnly,
					MinimumBlockSize:   uint32(blockSize),
					PreferredBlockSize: uint32(blockSize),
					MaximumBlockSize:   uint32(blockSize),
				}); err != nil {
				panic(err)
			}
		}()
	}
}
