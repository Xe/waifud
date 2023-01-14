package ninep

import (
	"encoding/binary"
	"errors"
	"io"
)

func writeString(w io.Writer, s string) error {
	b := []byte(s)
	if len(b) > 0xffff {
		return errors.New("string to write is too long")
	}
	if e := binary.Write(w, binary.LittleEndian, int16(len(b))); e != nil {
		return e
	}
	if e := binary.Write(w, binary.LittleEndian, b); e != nil {
		return e
	}
	return nil
}

func writeStringSlice(w io.Writer, ss []string) error {
	if err := writeUint16(w, uint16(len(ss))); err != nil {
		return err
	}
	for _, s := range ss {
		if err := writeString(w, s); err != nil {
			return err
		}
	}
	return nil
}

func writeQIDSlice(w io.Writer, qs []QID) error {
	if err := writeUint16(w, uint16(len(qs))); err != nil {
		return err
	}
	for _, q := range qs {
		if err := writeQID(w, q); err != nil {
			return err
		}
	}
	return nil
}

func writeByteSlice(w io.Writer, bs []byte) error {
	if err := writeUint32(w, uint32(len(bs))); err != nil {
		return err
	}
	n, err := w.Write(bs)
	if err != nil {
		return err
	}
	if n < len(bs) {
		// TODO: Repeat byte slice write instead.
		return errors.New("short write")
	}
	return nil
}

func writeQID(w io.Writer, q QID) error {
	// QID struct is laid out to serialize correctly.
	return binary.Write(w, binary.LittleEndian, q)
}

func writeUint8(w io.Writer, v uint8) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeUint16(w io.Writer, v uint16) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeUint32(w io.Writer, v uint32) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeUint64(w io.Writer, v uint64) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func stringSliceSize(ss []string) (size uint32) {
	size += 2 // number of strings in slice
	for _, s := range ss {
		size += 2 + uint32(len(s))
	}
	return size
}
