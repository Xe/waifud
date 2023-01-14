package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/qid"
	"io"
	"os"
)

//OpenMode is used to signalize the mode a file should be opened with in a TOpenMessage.
type OpenMode uint8

const (
	// Read opens a file for reading.
	Read OpenMode = 0

	// Write opens a file for writing
	Write OpenMode = 1

	// ReadWrite opens a file for random access.
	ReadWrite OpenMode = 2

	// Exec opens a file for reading but also checks execute permission.
	Exec OpenMode = 3

	// Trunc is to be OR'ed in (except for OpenExec), truncates file before opening.
	Trunc OpenMode = 16

	// CloseExec is to be OR'ed in, closes on execution.
	CloseExec OpenMode = 32

	// Close is to be OR'ed in, remove on closing.
	Close OpenMode = 64
)

// OsMode convertes a OpenMode to os.Open mode.
//
// Converted modes: Read, Write, ReadWrite, Trunc
func OsMode(mode OpenMode) int {
	flg := 0
	if mode&Read == Read {
		flg |= os.O_RDONLY
	}
	if mode&Write == Write {
		flg |= os.O_WRONLY
	}
	if mode&ReadWrite == ReadWrite {
		flg |= os.O_RDWR
	}
	if mode&Trunc == Trunc {
		flg |= os.O_TRUNC
	}

	return flg
}

// NeinMode convertes a os.Open mode to OpenMode.
//
// Converted modes: os.O_RDONLY, os.O_WRONLY, os.O_RDWR, os.O_TRUNC
func NeinMode(osMode int) OpenMode {
	var mode OpenMode
	if osMode&os.O_RDONLY == os.O_RDONLY {
		mode |= Read
	}
	if osMode&os.O_WRONLY == os.O_WRONLY {
		mode |= Write
	}
	if osMode&os.O_RDWR == os.O_RDWR {
		mode |= ReadWrite
	}
	if osMode&os.O_TRUNC == os.O_TRUNC {
		mode |= Trunc
	}

	return mode
}

/*TOpen is a 9P open request message.

Open requests ask the file server to prepare a fid to be used
with read and write requests, checking the permissions before.

Mode determines the type of IO and can be one of the constants
OpenRead, OpenWrite, OpenReadWrite, OpenExec, additionally
OpenTrunc, OpenCloseExec and OpenClose or'ed into Mode.

See also: http://man.cat-v.org/plan_9/5/open
*/
type TOpen struct {
	Fid  fid.Fid
	Mode OpenMode
}

func (m *TOpen) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint8Encode(w, uint8(m.Mode))
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}

func (m *TOpen) decode(r io.Reader) (int64, error) {
	f, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	mod, n2, err := basic.Uint8Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	m.Fid = fid.Fid(f)
	m.Mode = OpenMode(mod)

	return n1 + n2, nil
}

/*ROpen is the reply to a open request.

Qid is the servers idea of the opened file accessed.

Iounit may be zero. If not, it's the number of bytes which can
be read or written in a single call.

See also: http://man.cat-v.org/plan_9/5/open
*/
type ROpen struct {
	Qid    qid.Qid
	Iounit uint32
}

func (m *ROpen) encode(w io.Writer) (int64, error) {
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

func (m *ROpen) decode(r io.Reader) (int64, error) {
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
