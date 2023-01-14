package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestRErrorEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RError{Ename: "deadbeef"}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RError{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Log(a, b)
		t.Fail()
	}
}
