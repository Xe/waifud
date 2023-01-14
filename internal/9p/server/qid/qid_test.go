package qid // import "go.rbn.im/neinp/qid"

import (
	"bytes"
	"testing"
)

func TestQidEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &Qid{1, 2, 3}
	_, err := a.Encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &Qid{}
	_, err = b.Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
