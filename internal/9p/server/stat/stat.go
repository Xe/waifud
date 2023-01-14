//Package stat provides functionality to handle 9p stat values.
package stat // import "go.rbn.im/neinp/stat"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/qid"
	"bytes"
	"io"
	"time"
)

/*Stat contains file attributes.

Typ and Dev are usually unused. Qid contains the Qid of this file.
Mode is a combination of UNIX-permission bits for owner, group and others.
Atime contains the time of last access, Mtime the time of last modification.
Length is the file length in bytes, directories have a conventional
length of 0. Name is the file name and must be "/" if the file is the root
directory of the server. Uid is the owners name (not UNIX uid), group is the
groups name. Muid is the name of the last user who modified the file.

http://man.cat-v.org/inferno/5/stat
*/
type Stat struct {
	//	Size   uint16
	Typ    uint16
	Dev    uint32
	Qid    qid.Qid
	Mode   Mode
	Atime  time.Time
	Mtime  time.Time
	Length uint64
	Name   string
	Uid    string
	Gid    string
	Muid   string
}

//IsDir returns true if Dir is set.
func (s *Stat) IsDir() bool {
	return s.Mode&Dir == Dir
}

//IsAppend returns true if Append is set.
func (s *Stat) IsAppend() bool {
	return s.Mode&Append == Append
}

//IsExcl returns true if Excl is set.
func (s *Stat) IsExcl() bool {
	return s.Mode&Excl == Excl
}

//IsMount returns true if Mount is set.
func (s *Stat) IsMount() bool {
	return s.Mode&Mount == Mount
}

//IsAuth returns true if Auth is set.
func (s *Stat) IsAuth() bool {
	return s.Mode&Auth == Auth
}

//IsTmp returns true if Tmp is set.
func (s *Stat) IsTmp() bool {
	return s.Mode&Tmp == Tmp
}

//IsSymlink returns true if Symlink is set.
func (s *Stat) IsSymlink() bool {
	return s.Mode&Symlink == Symlink
}

//IsDevice returns true if Device is set.
func (s *Stat) IsDevice() bool {
	return s.Mode&Device == Device
}

//IsNamedPipe returns true if NamedPipe is set.
func (s *Stat) IsNamedPipe() bool {
	return s.Mode&NamedPipe == NamedPipe
}

//IsSocket returns true if Socket is set.
func (s *Stat) IsSocket() bool {
	return s.Mode&Socket == Socket
}

//IsSetUid returns true if SetUid is set.
func (s *Stat) IsSetUid() bool {
	return s.Mode&SetUid == SetUid
}

//IsSetGid returns true if SetGid is set
func (s *Stat) IsSetGid() bool {
	return s.Mode&SetGid == SetGid
}

func (s *Stat) Encode(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	_, err := basic.Uint16Encode(&buf, s.Typ)
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint32Encode(&buf, s.Dev)
	if err != nil {
		return 0, err
	}

	_, err = s.Qid.Encode(&buf)
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint32Encode(&buf, uint32(s.Mode))
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint32Encode(&buf, uint32(s.Atime.Unix()))
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint32Encode(&buf, uint32(s.Mtime.Unix()))
	if err != nil {
		return 0, err
	}

	_, err = basic.Uint64Encode(&buf, s.Length)
	if err != nil {
		return 0, err
	}

	_, err = basic.StringEncode(&buf, s.Name)
	if err != nil {
		return 0, err
	}

	_, err = basic.StringEncode(&buf, s.Uid)
	if err != nil {
		return 0, err
	}

	_, err = basic.StringEncode(&buf, s.Gid)
	if err != nil {
		return 0, err
	}

	_, err = basic.StringEncode(&buf, s.Muid)
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
	size := uint16(buf.Len()) // we don't need to subtract 2 because size isn't in buf

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

func (s *Stat) Decode(r io.Reader) (int64, error) {
	size, n1, err := basic.Uint16Decode(r)
	if err != nil {
		return n1, err
	}

	b := make([]byte, size)
	n2, err := r.Read(b)
	if err != nil {
		return n1 + int64(n2), err
	}

	n := n1 + int64(n2)

	buf := bytes.NewBuffer(b)

	typ, _, err := basic.Uint16Decode(buf)
	if err != nil {
		return n, err
	}

	dev, _, err := basic.Uint32Decode(buf)
	if err != nil {
		return n, err
	}

	_, err = s.Qid.Decode(buf)
	if err != nil {
		return n, err
	}

	mod, _, err := basic.Uint32Decode(buf)
	if err != nil {
		return n, err
	}

	atime, _, err := basic.Uint32Decode(buf)
	if err != nil {
		return n, err
	}

	mtime, _, err := basic.Uint32Decode(buf)
	if err != nil {
		return n, err
	}

	length, _, err := basic.Uint64Decode(buf)
	if err != nil {
		return n, err
	}

	name, _, err := basic.StringDecode(buf)
	if err != nil {
		return n, err
	}

	uid, _, err := basic.StringDecode(buf)
	if err != nil {
		return n, err
	}

	gid, _, err := basic.StringDecode(buf)
	if err != nil {
		return n, err
	}

	muid, _, err := basic.StringDecode(buf)
	if err != nil {
		return n, err
	}

	s.Typ = typ
	s.Dev = dev
	s.Mode = Mode(mod)
	s.Atime = time.Unix(int64(atime), 0)
	s.Mtime = time.Unix(int64(mtime), 0)
	s.Length = length
	s.Name = name
	s.Uid = uid
	s.Gid = gid
	s.Muid = muid

	return n, nil
}
