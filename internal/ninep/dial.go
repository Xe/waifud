package ninep

import (
	"context"
	"fmt"
	"log"
	"net"
)

// nofid is the fid value used to indicate absence of a FID,
// e.g. to pass as afid when no authentication is required.
const nofid uint32 = ^uint32(0)

// notag is the tag value used in absence of a tag,
// e.g. during authentication
const notag uint16 = ^uint16(0)

func dialNet(service string) (net.Conn, error) {
	if service == "sources" {
		return net.Dial("tcp", "sources.9p.io:564")
	}
	return net.Dial("tcp", service)
}

func handshake(c net.Conn) (msize uint32, err error) {
	uname, aname := "user", ""
	wantVersion := "9P2000"
	var wantMsize uint32 = 8192
	rootFID := uint32(0) // TODO: Dynamically acquire FIDs somehow

	if err := writeTversion(c, notag, wantMsize, wantVersion); err != nil {
		return 0, err
	}
	msize, version, err := readRversion(c)
	if err != nil {
		return 0, fmt.Errorf("version(%q, %q): %w", wantMsize, wantVersion, err)
	}

	if wantMsize < msize {
		return 0, fmt.Errorf("server wanted too high msize of %v", msize)
	}
	if version != wantVersion {
		return 0, fmt.Errorf("mismatching version: %q != %q", version, wantVersion)
	}

	// Afid is nofid when the client doesn't want to authenticate.
	afid := nofid

	// XXX: Authentication step

	if err := writeTattach(c, 1, rootFID, afid, uname, aname); err != nil {
		return 0, err
	}
	_, err = readRattach(c)
	if err != nil {
		return 0, fmt.Errorf("attach(): %w", err)
	}
	return msize, nil
}

type dialOptions struct {
	concurrency uint16
}

type dialOpt func(*dialOptions)

func WithConcurrency(concurrency uint16) dialOpt {
	return func(c *dialOptions) {
		c.concurrency = concurrency
	}
}

func Dial(service string, opts ...dialOpt) (*FS, error) {
	cc, err := dial9pConn(service, opts...)
	if err != nil {
		return nil, err
	}
	return &FS{cc: cc}, nil
}

// dial9pConn establishes a 9p client connection and returns it.
func dial9pConn(service string, opts ...dialOpt) (*clientConn, error) {
	options := dialOptions{
		concurrency: 256,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Dial
	netConn, err := dialNet(service)
	if err != nil {
		return nil, err
	}

	// Handshake
	msize, err := handshake(netConn)
	if err != nil {
		netConn.Close()
		return nil, err
	}

	// Build client connection.
	cc := &clientConn{
		tags:       make(chan uint16, options.concurrency),
		conn:       netConn,
		reqReaders: make(map[uint16]callback),
		msize:      msize,
	}
	go func() {
		for i := uint16(0); i < options.concurrency; i++ {
			cc.tags <- i
		}
	}()

	go func() {
		err := cc.run(context.Background()) // TODO: Cancelation
		if err != nil {
			// TODO: How to report error correctly?
			log.Fatalf("9p client: run(): %v", err)
		}
	}()

	return cc, nil
}
