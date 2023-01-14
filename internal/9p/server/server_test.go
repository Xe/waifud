package neinp

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/Xe/waifud/internal/9p/server/message"
)

type MsizeFS struct {
	NopP2000
}

func (f *MsizeFS) Version(ctx context.Context, req message.TVersion) (message.RVersion, error) {
	return message.RVersion{Msize: req.Msize, Version: req.Version}, nil
}

func TestMsize(t *testing.T) {
	s := NewServer(&MsizeFS{})

	content := message.TVersion{Msize: 1234, Version: "9P2000"}
	msg, err := message.New(0x1, &content)
	if err != nil {
		t.Error(err)
	}

	s.process(msg)

	if s.msize != content.Msize {
		t.Logf("Want: %v Have: %v", content.Msize, s.msize)
		t.Fail()
	}
}

func TestRcv(t *testing.T) {
	buf := bytes.NewReader([]byte{0x1b, 0x0, 0x0, 0x0, 0x68, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x74, 0x65, 0x73, 0x74, 0x4, 0x0, 0x74, 0x65, 0x73, 0x74})

	done := make(chan struct{})

	s := &Server{}

	msgs, errs := s.rcv(done, buf)

	select {
	case x := <-msgs:
		y, ok := x.Content.(*message.TAttach)
		if !ok {
			t.Fail()
		}
		if y.Uname != "test" {
			t.Fail()
		}
		if y.Aname != "test" {
			t.Fail()
		}
	case err := <-errs:
		if err != io.EOF {
			t.Error(err)
		}
	}

}
