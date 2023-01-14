package client

import (
	"encoding/binary"
	"io"
)

func readString(r io.Reader, s *string) error {
	var sz uint16
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return err
	}
	buf := make([]byte, sz)
	if err := binary.Read(r, binary.LittleEndian, &buf); err != nil {
		return err
	}
	*s = string(buf)
	return nil
}

func readStringSlice(r io.Reader, ss *[]string) error {
	var size int16
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}
	*ss = make([]string, 0, size)
	for i := int16(0); i < size; i++ {
		var s string
		if err := readString(r, &s); err != nil {
			return err
		}
		*ss = append(*ss, s)
	}
	return nil
}

func readQIDSlice(r io.Reader, qs *[]QID) error {
	var size uint16
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}
	*qs = make([]QID, 0, size)
	for i := uint16(0); i < size; i++ {
		var q QID
		if err := readQID(r, &q); err != nil {
			return err
		}
		*qs = append(*qs, q)
	}
	return nil
}

func readQID(r io.Reader, q *QID) error {
	return binary.Read(r, binary.LittleEndian, q)
}

// Note: This *populates* a byte slice passed in from the outside.
func readAndFillByteSlice(r io.Reader, bs []byte) (uint32, error) {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return 0, err
	}
	n := size
	if n > uint32(len(bs)) {
		n = uint32(len(bs))
	}
	bs = bs[:n]
	if _, err := io.ReadFull(r, bs); err != nil {
		return 0, err
	}
	if err := skip(r, int(size-n)); err != nil {
		return 0, err
	}
	return n, nil
}

func readUint8(r io.Reader, out *uint8) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint16(r io.Reader, out *uint16) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint32(r io.Reader, out *uint32) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func readUint64(r io.Reader, out *uint64) error {
	return binary.Read(r, binary.LittleEndian, out)
}

func skip(r io.Reader, n int) error {
	if s, ok := r.(io.ReadSeeker); ok {
		_, err := s.Seek(int64(n), io.SeekCurrent)
		return err
	}
	var buffer [1024]byte
	buf := buffer[:]
	for n > 0 {
		if n < len(buf) {
			buf = buf[:n]
		}
		nn, err := r.Read(buf)
		if err != nil {
			return err
		}
		n -= nn
	}
	return nil
}
