package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/stat"
	"bytes"
	"io"
)

/*TStat requests information about file.

The file for which information is requested is identified by Fid.

See also: http://man.cat-v.org/inferno/5/stat
*/
type TStat struct {
	Fid fid.Fid
}

func (m *TStat) encode(w io.Writer) (int64, error) {
	n, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *TStat) decode(r io.Reader) (int64, error) {
	f, n, err := basic.Uint32Decode(r)
	if err != nil {
		return n, err
	}

	m.Fid = fid.Fid(f)
	return n, nil
}

/*RStat contains a machine-independent directory entry, Stat.

See also: http://man.cat-v.org/inferno/5/stat
*/
type RStat struct {
	Stat stat.Stat
}

func (m *RStat) encode(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := m.Stat.Encode(&buf)
	if err != nil {
		return 0, err
	}

	/* From plan9 manual:
	BUGS
	To make the contents of a directory, such as returned by
	read(9P), easy to parse, each directory entry begins with a
	size field.  For consistency, the entries in Twstat and
	Rstat messages also contain their size, which means the size
	appears twice.  For example, the Rstat message is formatted
	as ``(4+1+2+2+n)[4] Rstat tag[2] n[2] (n-2)[2] type[2]
	dev[4]...,'' where n is the value returned by convD2M.
	*/

	size := uint16(buf.Len())
	n1, err := basic.Uint16Encode(w, size)
	if err != nil {
		return n1, err
	}

	n2, err := buf.WriteTo(w)
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}

func (m *RStat) decode(r io.Reader) (int64, error) {
	// read the size but don't store it. see Rstat.encode
	_, n1, err := basic.Uint16Decode(r)
	if err != nil {
		return n1, err
	}

	n2, err := m.Stat.Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}
