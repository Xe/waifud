package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/qid"
	"bytes"
	"strings"
	"testing"
)

func TestTWalkEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TWalk{Fid: 1, Newfid: 2, Wname: []string{"dead", "beef"}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TWalk{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if a.Fid != b.Fid || a.Newfid != b.Newfid {
		t.Log(a, b)
		t.Fail()
	}

	if len(a.Wname) != len(b.Wname) {
		t.Fail()
	}

	if strings.Join(a.Wname, "") != strings.Join(b.Wname, "") {
		t.Fail()
	}
}

func TestRWalkEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RWalk{Wqid: []qid.Qid{qid.Qid{1, 2, 3}, qid.Qid{4, 5, 6}}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RWalk{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if len(a.Wqid) != len(b.Wqid) {
		t.Logf("wrong wqid len, a: %v, b: %v", len(a.Wqid), len(b.Wqid))
		t.Fail()
	} else {
		for i, v := range a.Wqid {
			if b.Wqid[i] != v {
				t.Logf("wqids differ, a.Wqid[i]: %v, b.Wqid[i]: %v", v, b.Wqid[i])
				t.Fail()
			}
		}
	}
}
