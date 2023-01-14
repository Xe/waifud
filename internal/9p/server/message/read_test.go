package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTReadEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TRead{Fid: 1, Offset: 2, Count: 3}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TRead{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRReadEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RRead{Count: 4, Data: []byte{0xDE, 0xAD, 0xBE, 0xEF}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RRead{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if a.Count != b.Count || !bytes.Equal(a.Data, b.Data) {
		t.Log(a, b)
		t.Fail()
	}
}
