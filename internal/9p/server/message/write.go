package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"io"
)

/*TWrite requests data to be written to a file.

Write Count bytes of Data at Offset bytes from the beginning
of the file. The file which is represented by Fid must be
opened for writing before this request, when opened in append
mode, the write will be to the end of the file, ignoring Offset.
Writes to directories are forbidden.

See also: http://man.cat-v.org/inferno/5/read
*/
type TWrite struct {
	Fid    fid.Fid
	Offset uint64
	Count  uint32
	Data   []byte
}

func (m *TWrite) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint64Encode(w, m.Offset)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := basic.Uint32Encode(w, m.Count)
	if err != nil {
		return n1 + n2 + n3, err
	}

	n4, err := w.Write(m.Data)
	if err != nil {
		return n1 + n2 + n3 + int64(n4), err
	}

	return n1 + n2 + n3 + int64(n4), nil
}

func (m *TWrite) decode(r io.Reader) (int64, error) {
	f, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	offset, n2, err := basic.Uint64Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	count, n3, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	m.Fid = fid.Fid(f)
	m.Offset = offset
	m.Count = count

	m.Data = make([]byte, count)

	n4, err := r.Read(m.Data)
	if err != nil {
		return n1 + n2 + n3 + int64(n4), err
	}

	return n1 + n2 + n3 + int64(n4), nil
}

/*RWrite is the reply to a TWrite.

Count is the number of bytes written. It usually should be
the same as requested.

See also: http://man.cat-v.org/inferno/5/read
*/
type RWrite struct {
	Count uint32
}

func (m *RWrite) encode(w io.Writer) (int64, error) {
	n, err := basic.Uint32Encode(w, m.Count)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *RWrite) decode(r io.Reader) (int64, error) {
	count, n, err := basic.Uint32Decode(r)
	if err != nil {
		return n, err
	}

	m.Count = count

	return n, nil
}
