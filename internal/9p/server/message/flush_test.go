package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTFlushEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TFlush{Oldtag: 1}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TFlush{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRFlushEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RFlush{}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RFlush{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
