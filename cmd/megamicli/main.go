package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/Xe/waifud/internal/9p/client"
)

var (
	addr = flag.String("addr", "127.0.0.1:564", "location of megamid on tailnet")
)

func main() {
	flag.Parse()

	fs, err := client.Dial(*addr)
	if err != nil {
		log.Fatal(err)
	}

	fin, err := fs.Open("foo")
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	io.Copy(os.Stdout, fin)
}
