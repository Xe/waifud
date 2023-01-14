package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTRemoveEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TRemove{Fid: 1}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TRemove{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRRemoveEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RRemove{}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RRemove{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
