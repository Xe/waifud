package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/qid"
	"go.rbn.im/neinp/stat"
	"io"
)

/*TCreate asks the file server to create a new file.

The file is to be created in the directory represented by Fid, with the supplied Name.
Write permission is required on the directory. The owner will be implied by the
file systems used. The group is the same as the directories.
Permissions for a file will be set to

	Perm & 0666 | DPerm & 0666

or

	Perm & 0777 | DPerm & 0777

if a directory (DPerm being the permissions of the parent directory).

The created file is opened with Mode, and Fid will represent the new file.
Mode is not checked to fulfill the permissions of Perm.

Directories are created by setting the DirModeDir bit
in Perm.

Names "." and ".." are forbidden.

Fid must not be in use already.

Creating a file with a name already in use, will truncate the existing file.

See also: http://man.cat-v.org/plan_9/5/open
*/
type TCreate struct {
	Fid  uint32
	Name string
	Perm stat.Mode
	Mode OpenMode
}

func (m *TCreate) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, m.Fid)
	if err != nil {
		return n1, err
	}

	n2, err := basic.StringEncode(w, m.Name)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := basic.Uint32Encode(w, uint32(m.Perm))
	if err != nil {
		return n1 + n2 + n3, err
	}

	n4, err := basic.Uint8Encode(w, uint8(m.Mode))
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	return n1 + n2 + n3 + n4, nil
}

func (m *TCreate) decode(r io.Reader) (int64, error) {
	fid, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	name, n2, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2, err
	}

	perm, n3, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	mod, n4, err := basic.Uint8Decode(r)
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	m.Fid = fid
	m.Name = name
	m.Perm = stat.Mode(perm)
	m.Mode = OpenMode(mod)

	return n1 + n2 + n3 + n4, err
}

/*RCreate is the answer for a successful create.

Qid is the new file as seen by the server.

Iounit may be zero. If not, it is the number of bytes guaranteed
to succeed to be read or written in one message.

See also: http://man.cat-v.org/plan_9/5/open
*/
type RCreate struct {
	Qid    qid.Qid
	Iounit uint32
}

func (m *RCreate) encode(w io.Writer) (int64, error) {
	n1, err := m.Qid.Encode(w)
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint32Encode(w, m.Iounit)
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}

func (m *RCreate) decode(r io.Reader) (int64, error) {
	n1, err := m.Qid.Decode(r)
	if err != nil {
		return n1, err
	}

	iounit, n2, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	m.Iounit = iounit

	return n1 + n2, nil
}
