package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/stat"
	"bytes"
	"io"
)

/*TWstat requests to change some of the file status information.

Name can be changed by anyone with write permissions for the parent directory
and if the name isn't already in use in the parent directory.

The length of the file can be changed by anyone with write
permissions and will reflect in the actual length of the file.
The length of directories can't be set to anything than zero.

Servers may refuse length changes.

Mode and mtime may be changed by the owner or group.

All permission and mode bits can be changed, except for
the directory bit.

The gid can be changed by the owner if he is member of the
new group or by the group leader of both the current and new group.

The other data can't be changed using wstat. It is illegal to change
the owner of a file.

Changes are performed atomic, either all changes are applied with the
call succeeding, or none.

Modification of properties can be avoided using zero values:
Strings of length zero and the maximum value for integers.

If every property is set to such a zero value, it is to be handled
as a request to sync the file to stable storage.

See also: http://man.cat-v.org/inferno/5/stat
*/
type TWstat struct {
	Fid  fid.Fid
	Stat stat.Stat
}

func (m *TWstat) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n1, err
	}

	var buf bytes.Buffer

	_, err = m.Stat.Encode(&buf)
	if err != nil {
		return n1, err
	}

	size := uint16(buf.Len())
	n2, err := basic.Uint16Encode(w, size)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := m.Stat.Encode(w)
	if err != nil {
		return n1 + n2 + n3, err
	}

	return n1 + n2 + n3, err
}

func (m *TWstat) decode(r io.Reader) (int64, error) {
	f, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	_, n2, err := basic.Uint16Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := m.Stat.Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	m.Fid = fid.Fid(f)

	return n1 + n2 + n3, nil
}

/*RWstat is the reply for a successfull TWstat request.

See also: http://man.cat-v.org/inferno/5/stat
*/
type RWstat struct{}

func (m *RWstat) encode(w io.Writer) (int64, error) {
	return 0, nil
}

func (m *RWstat) decode(r io.Reader) (int64, error) {
	return 0, nil
}
