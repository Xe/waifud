package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Xe/waifud/internal/ninep"
)

var (
	addr    = flag.String("addr", "localhost:8080", "Address to serve HTTP on")
	service = flag.String("service", "sources", "9p service to connect to")
)

func main() {
	fsys, err := ninep.Dial(*service)
	if err != nil {
		log.Fatalf("ninep.Dial(%q): %v", *service, err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe(%q, nil): %v", *addr, err)
	}
}
