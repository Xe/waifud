package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"regexp"
	"sync"
)

var (
	ip2mac     = map[string]string{}
	ip2macLock sync.RWMutex

	dnsmasqParser = regexp.MustCompile(`DHCPACK\(\w+\) ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+) ([0-9a-f][0-9a-f]:[0-9a-f][0-9a-f]:[0-9a-f][0-9a-f]:[0-9a-f][0-9a-f]:[0-9a-f][0-9a-f]:[0-9a-f][0-9a-f]) (\w+)`)
)

func streamSyslog(ctx context.Context, program string) error {
	journalctlPath, err := exec.LookPath("journalctl")
	if err != nil {
		return fmt.Errorf("can't find journalctl: %w", err)
	}

	rdr, wtr := net.Pipe()

	cmd := exec.CommandContext(ctx, journalctlPath, "-f", "_COMM="+program, "--output", "cat")
	cmd.Stdout = wtr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("can't run command: %w", err)
	}

	go readAndParseDHCP(ctx, rdr)

	return nil
}

func readAndParseDHCP(ctx context.Context, r io.ReadCloser) error {
	scnr := bufio.NewScanner(r)
	defer r.Close()

	for scnr.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line := scnr.Text()
		parts := dnsmasqParser.FindStringSubmatch(line)
		if parts == nil {
			continue
		}
		ip := parts[1]
		macAddr := parts[2]

		ip2macLock.Lock()
		ip2mac[ip] = macAddr
		ip2macLock.Unlock()
	}

	return nil
}
