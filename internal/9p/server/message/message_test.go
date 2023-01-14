package message // import "go.rbn.im/neinp/message"

import (
	"bytes"
	"testing"
)

var contentTypeTests = []struct {
	c   Content
	typ messageType
}{
	{c: &TVersion{}, typ: tversion},
	{c: &RVersion{}, typ: rversion},
	{c: &TAuth{}, typ: tauth},
	{c: &RAuth{}, typ: rauth},
	{c: &TAttach{}, typ: tattach},
	{c: &RAttach{}, typ: rattach},
	{c: &RError{}, typ: rerror},
	{c: &TFlush{}, typ: tflush},
	{c: &RFlush{}, typ: rflush},
	{c: &TWalk{}, typ: twalk},
	{c: &RWalk{}, typ: rwalk},
	{c: &TOpen{}, typ: topen},
	{c: &ROpen{}, typ: ropen},
	{c: &TCreate{}, typ: tcreate},
	{c: &RCreate{}, typ: rcreate},
	{c: &TRead{}, typ: tread},
	{c: &RRead{}, typ: rread},
	{c: &TWrite{}, typ: twrite},
	{c: &RWrite{}, typ: rwrite},
	{c: &TClunk{}, typ: tclunk},
	{c: &RClunk{}, typ: rclunk},
	{c: &TRemove{}, typ: tremove},
	{c: &RRemove{}, typ: rremove},
	{c: &TStat{}, typ: tstat},
	{c: &RStat{}, typ: rstat},
	{c: &TWstat{}, typ: twstat},
	{c: &RWstat{}, typ: rwstat},
}

func TestContentType(t *testing.T) {
	for _, v := range contentTypeTests {
		x, err := contentType(v.c)
		if err != nil {
			t.Error(err)
		}

		if v.typ != x {
			t.Logf("%T %v %v", v.c, v.typ, x)
			t.Fail()
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	mIn, err := New(0x1, &TVersion{Version: "9P2000.test", Msize: 1234})
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer

	n1, err := mIn.Encode(&buf)
	if err != nil {
		t.Error(err)
	}

	var mOut Message

	n2, err := mOut.Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if n1 != n2 {
		t.Fail()
	}
}
