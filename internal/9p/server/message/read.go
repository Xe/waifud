package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"io"
)

/*TRead requests data from a file.

The file is identified by Fid, which must be opened before reading.
Count bytes are read from the file, starting at Offset bytes after the
beginning of the file.

The count Field in the reply indicates the number of bytes
returned.  This may be less than the requested amount.  If
the Offset field is greater than or equal to the number of
bytes in the file, a Count of zero will be returned.

For directories, the read request message must have
Offset equal to zero or the value of Offset in the previous
read on the directory, plus the number of bytes returned in
the previous read. In other words, seeking other than to
the beginning is illegal in a directory.

See also: http://man.cat-v.org/inferno/5/read
*/
type TRead struct {
	Fid    fid.Fid
	Offset uint64
	Count  uint32
}

func (m *TRead) encode(w io.Writer) (int64, error) {
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

	return n1 + n2 + n3, nil
}

func (m *TRead) decode(r io.Reader) (int64, error) {
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

	return n1 + n2 + n3, nil
}

/*RRead is the reply to a TRead.

The data read is stored in Data, Count being the length
of Data in bytes.

For directories, read returns an integral number of directory
entries exactly as in RStat, one for each member of the directory.
To help with this StatReader can be used.

See also: http://man.cat-v.org/inferno/5/read
*/
type RRead struct {
	Count uint32
	Data  []byte
}

func (m *RRead) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, m.Count)
	if err != nil {
		return n1, err
	}

	n2, err := w.Write(m.Data)
	if err != nil {
		return n1 + int64(n2), err
	}

	return n1 + int64(n2), nil
}

func (m *RRead) decode(r io.Reader) (int64, error) {
	count, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	m.Count = count
	m.Data = make([]byte, count)

	n2, err := r.Read(m.Data)
	if err != nil {
		return n1 + int64(n2), err
	}

	return n1 + int64(n2), nil
}
