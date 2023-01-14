package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"bytes"
	"testing"
)

func TestTOpenEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TOpen{Fid: 1, Mode: 2}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TOpen{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestROpenEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &ROpen{Qid: qid.Qid{Type: 1, Version: 2, Path: 3}, Iounit: 4}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &ROpen{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
