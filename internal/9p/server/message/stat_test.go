package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"go.rbn.im/neinp/stat"
	"bytes"
	"testing"
	"time"
)

func TestStatEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &stat.Stat{Typ: 1, Dev: 2, Qid: qid.Qid{3, 4, 5}, Mode: 6, Atime: time.Unix(1234567, 0), Mtime: time.Unix(1234567, 0), Length: 6, Name: "deadbeef", Uid: "foo", Gid: "bar", Muid: "baz"}
	_, err := a.Encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &stat.Stat{}
	_, err = b.Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestTStatEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TStat{Fid: 1}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TStat{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}

func TestRStatEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RStat{Stat: stat.Stat{Typ: 1, Dev: 2, Qid: qid.Qid{3, 4, 5}, Mode: 6, Atime: time.Unix(1234567, 0), Mtime: time.Unix(1234567, 0), Length: 6, Name: "deadbeef", Uid: "foo", Gid: "bar", Muid: "baz"}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RStat{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
