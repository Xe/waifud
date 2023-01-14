package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	ninep "github.com/Xe/waifud/internal/9p/client"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage\n")
	fmt.Fprintf(flag.CommandLine.Output(), "     %s CMD PATH\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "e.g. %s cat sources/plan9/NOTICE\n\n", os.Args[0])
	flag.PrintDefaults()
}

func parsePositionalArgs() (cmd string, service string, path string) {
	if len(flag.Args()) != 2 {
		log.Fatal("Could not parse positional args; want [cmd] [path]")
	}
	cmd = flag.Args()[0]
	arg := flag.Args()[1]
	service, path, _ = strings.Cut(arg, "/")
	return
}

func formatStat(stat fs.FileInfo) string {
	sys := stat.Sys().(ninep.Stat)
	return fmt.Sprintf("%s %8d %8s %8s %s",
		stat.Mode().String(), stat.Size(), sys.UID, sys.GID, stat.Name())
}

func main() {
	// For better RPC latency debugging, log microseconds.
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.Usage = usage
	flag.Parse()
	cmd, service, path := parsePositionalArgs()

	fsys, err := ninep.Dial(service)
	if err != nil {
		log.Fatalf("Dial(%q): %v", service, err)
	}

	switch cmd {
	case "cat":
		r, err := fsys.Open(path)
		if err != nil {
			log.Fatalf("Open: %v", err)
		}
		defer r.Close()
		buf, err := io.ReadAll(r)
		if err != nil {
			log.Fatalf("Read: %v", err)
		}
		os.Stdout.Write(buf)

	case "stat":
		stat, err := fs.Stat(fsys, path)
		if err != nil {
			log.Fatalf("Stat: %v", err)
		}
		fmt.Println(formatStat(stat))

	case "ls":
		entries, err := fs.ReadDir(fsys, path)
		if err != nil {
			log.Fatalf("ReadDir: %v", err)
		}
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				log.Fatalf("Info(): %v", err)
			}
			fmt.Println(formatStat(info))
		}
	}
}
