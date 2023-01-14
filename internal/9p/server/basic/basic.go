//Package basic handles en-/decoding of basic types to 9p wire format.
package basic // import "go.rbn.im/neinp/basic"

import (
	"io"
)

//Uint8Decode reads a single byte.
//
//Returns the decoded byte, the count of read bytes and a possible error.
func Uint8Decode(r io.Reader) (uint8, int64, error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if n < 1 || err != nil {
		return 0, int64(n), err
	}
	return uint8(buf[0]), int64(n), nil
}

//Uint8Encode writes a single byte.
//
//Returns the count of written bytes and a possible error.
func Uint8Encode(w io.Writer, x uint8) (int64, error) {
	buf := []byte{x}
	n, err := w.Write(buf)
	if n < 1 || err != nil {
		return int64(n), err
	}
	return int64(n), nil
}

//Uint16Decode reads a 16 bit value from an io.Reader.
//
//Returns the decoded uint16 value, the count of read bytes and a possible error.
func Uint16Decode(r io.Reader) (uint16, int64, error) {
	buf := make([]byte, 2)
	n, err := r.Read(buf)
	if n < 2 || err != nil {
		return 0, int64(n), err
	}
	x := uint16(buf[0]) | uint16(buf[1])<<8
	return x, int64(n), nil
}

//Uint16Encode writes a uint16 value.
//
//Returns the count of written bytes and a possible error.
func Uint16Encode(w io.Writer, x uint16) (int64, error) {
	buf := []byte{uint8(x), uint8(x >> 8)}
	n, err := w.Write(buf)
	if n < 2 || err != nil {
		return int64(n), err
	}
	return int64(n), nil
}

//Uint32Decode reads a 32 bit value from an io.Reader.
//
//Returns the decoded uint32 value, the count of read bytes and a possible error.
func Uint32Decode(r io.Reader) (uint32, int64, error) {
	buf := make([]byte, 4)
	n, err := r.Read(buf)
	if n < 4 || err != nil {
		return 0, int64(n), err
	}
	x := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
	return x, int64(n), nil
}

//Uint32Encode writes a uint32 value.
//
//Returns the count of written bytes and a possible error.
func Uint32Encode(w io.Writer, x uint32) (int64, error) {
	buf := []byte{uint8(x), uint8(x >> 8), uint8(x >> 16), uint8(x >> 24)}
	n, err := w.Write(buf)
	if n < 4 || err != nil {
		return int64(n), err
	}
	return int64(n), nil
}

//Uint64Decode reads a 64 bit value from an io.Reader.
//
//Returns the decoded uint64 value, the count of read bytes and a possible error.
func Uint64Decode(r io.Reader) (uint64, int64, error) {
	buf := make([]byte, 8)
	n, err := r.Read(buf)
	if n < 8 || err != nil {
		return 0, int64(n), err
	}
	x := uint64(buf[0]) | uint64(buf[1])<<8 | uint64(buf[2])<<16 | uint64(buf[3])<<24 | uint64(buf[4])<<32 | uint64(buf[5])<<40 | uint64(buf[6])<<48 | uint64(buf[7])<<56
	return x, int64(n), nil
}

//Uint64Encode writes a uint64 value.
//
//Returns the count of written bytes and a possible error.
func Uint64Encode(w io.Writer, x uint64) (int64, error) {
	buf := []byte{uint8(x), uint8(x >> 8), uint8(x >> 16), uint8(x >> 24), uint8(x >> 32), uint8(x >> 40), uint8(x >> 48), uint8(x >> 56)}
	n, err := w.Write(buf)
	if n < 8 || err != nil {
		return int64(n), err
	}
	return int64(n), nil
}

//BytesDecode reads a variable number of bytes into a byte slice.
//
//Returns a byte slice, the count of read bytes and a possible error.
func BytesDecode(r io.Reader) ([]byte, int64, error) {
	size, n1, err := Uint16Decode(r)
	if err != nil {
		return []byte{}, n1, err
	}

	buf := make([]byte, size)
	n2, err := r.Read(buf)
	if err != nil {
		return []byte{}, n1 + int64(n2), err
	}

	return buf, n1 + int64(n2), nil
}

//BytesEncode writes a byte slice.
//
//Returns the count of written bytes and a possible error.
func BytesEncode(w io.Writer, b []byte) (int64, error) {
	size := uint16(len(b))
	n1, err := Uint16Encode(w, size)
	if err != nil {
		return n1, err
	}

	n2, err := w.Write(b)
	if err != nil {
		return n1 + int64(n2), err
	}

	return n1 + int64(n2), nil
}

//StringDecode reads a string.
//
//Returns a string, the count of read bytes and a possible error.
func StringDecode(r io.Reader) (string, int64, error) {
	b, n1, err := BytesDecode(r)
	if err != nil {
		return "", n1, err
	}

	return string(b), n1, nil
}

//StringEncode writes a string.
//
//Returns the count of written bytes and a possible error.
func StringEncode(w io.Writer, s string) (int64, error) {
	return BytesEncode(w, []byte(s))
}
