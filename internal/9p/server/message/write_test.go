package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

func TestTWriteEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &TWrite{Fid: 1, Offset: 2, Count: 4, Data: []byte{0xDE, 0xAD, 0xBE, 0xEF}}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &TWrite{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if a.Fid != b.Fid || a.Offset != b.Offset || a.Count != b.Count {
		t.Log(a, b)
		t.Fail()
	}

	if !bytes.Equal(a.Data, b.Data) {
		t.Log(a, b)
		t.Log(a.Data, b.Data)
		t.Fail()
	}
}

func TestRWriteEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	a := &RWrite{Count: 1}
	_, err := a.encode(&buf)
	if err != nil {
		t.Error(err)
	}

	b := &RWrite{}
	_, err = b.decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if *a != *b {
		t.Fail()
	}
}
