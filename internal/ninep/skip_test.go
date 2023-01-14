package ninep

import (
	"io"
	"strings"
	"testing"
)

func repeat(s string, n int) (r string) {
	for ; n > 0; n-- {
		r += s
	}
	return
}

func TestSkipSimple(t *testing.T) {
	s := strings.NewReader("foobar")

	if err := skip(s, 3); err != nil {
		t.Errorf("skip: %v", err)
	}

	if s.Len() != 3 {
		t.Errorf("expected remaining length 3, got %d", s.Len())
	}
}

type forwardReader struct {
	io.Reader
}

func TestSkipLongWithoutSeek(t *testing.T) {
	s := strings.NewReader(repeat("x", 5000))
	// Wrap it in forwardReader, so we're sure we don't
	// accidentally expose more than io.Reader.
	r := forwardReader{s}

	if err := skip(r, 3000); err != nil {
		t.Errorf("skip: %v", err)
	}

	want := 2000
	if s.Len() != want {
		t.Errorf("expected remaining length %d, got %d", want, s.Len())
	}
}
