package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"bytes"
	"testing"
)

func TestTAttachEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TAttach{Fid: 1, Afid: 2, Uname: "3", Aname: "4"}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TAttach{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRAttachEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RAttach{Qid: qid.Qid{Type: 1, Version: 2, Path: 3}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RAttach{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
