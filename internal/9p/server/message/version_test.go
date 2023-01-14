package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTVersionEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TVersion{Msize: 1, Version: "deadbeef"}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TVersion{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRVersionEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RVersion{Msize: 1, Version: "deadbeef"}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RVersion{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
