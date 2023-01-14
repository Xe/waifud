package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTClunkEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TClunk{Fid: 1}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TClunk{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRClunkEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RClunk{}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RClunk{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
