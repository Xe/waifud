package basic // import "go.rbn.im/neinp/basic"

import (
	"bytes"
	"testing"
)

func TestUint8EncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := Uint8Encode(&buf, 1<<8-1)
	if err != nil {
		t.Error(err)
	}

	x, _, err := Uint8Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if x != 1<<8-1 {
		t.Fail()
	}
}

func TestUint16EncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := Uint16Encode(&buf, 1<<16-1)
	if err != nil {
		t.Error(err)
	}

	x, _, err := Uint16Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if x != 1<<16-1 {
		t.Fail()
	}
}

func TestUint32EncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := Uint32Encode(&buf, 1<<32-1)
	if err != nil {
		t.Error(err)
	}

	x, _, err := Uint32Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if x != 1<<32-1 {
		t.Fail()
	}
}

func TestUint64EncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := Uint64Encode(&buf, 1<<64-1)
	if err != nil {
		t.Error(err)
	}

	x, _, err := Uint64Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if x != 1<<64-1 {
		t.Fail()
	}
}

func TestBytesEncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := BytesEncode(&buf, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		t.Error(err)
	}

	x, _, err := BytesDecode(&buf)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(x, []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
		t.Fail()
	}
}

func TestStringEncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	_, err := StringEncode(&buf, "deadbeef")
	if err != nil {
		t.Error(err)
	}

	x, _, err := StringDecode(&buf)
	if err != nil {
		t.Error(err)
	}

	if x != "deadbeef" {
		t.Fail()
	}
}
