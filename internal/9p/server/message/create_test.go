package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"bytes"
	"testing"
)

func TestTCreateEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TCreate{Fid: 1, Name: "a", Perm: 2, Mode: 3}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TCreate{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRCreateEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RCreate{Qid: qid.Qid{1, 2, 3}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RCreate{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
