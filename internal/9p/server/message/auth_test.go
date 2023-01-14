package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"bytes"
	"testing"
)

func TestTAuthEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TAuth{Afid: 1, Uname: "a", Aname: "b"}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TAuth{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRAuthEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RAuth{Aqid: qid.Qid{Type: 1, Version: 2, Path: 3}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RAuth{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
