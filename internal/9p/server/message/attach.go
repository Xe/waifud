package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/qid"
	"io"
)

/*TAttach establishes a new connection.

Fid will be mapped to the root of the file system.
Afid specifies a fid gained from a previous TAuth.
Uname identifies the user accessing the file system, Aname
selecting the file tree to access.

See also: http://man.cat-v.org/plan_9/5/attach
*/
type TAttach struct {
	Fid   fid.Fid
	Afid  fid.Fid
	Uname string
	Aname string
}

func (m *TAttach) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint32Encode(w, uint32(m.Afid))
	if err != nil {
		return n1 + n2, err
	}

	n3, err := basic.StringEncode(w, m.Uname)
	if err != nil {
		return n1 + n2 + n3, err
	}

	n4, err := basic.StringEncode(w, m.Aname)
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	return n1 + n2 + n3 + n4, nil
}

func (m *TAttach) decode(r io.Reader) (int64, error) {
	f, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	af, n2, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	uname, n3, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	aname, n4, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	m.Fid = fid.Fid(f)
	m.Afid = fid.Fid(af)
	m.Uname = uname
	m.Aname = aname

	return n1 + n2 + n3 + n4, err
}

/*RAttach is the answer to an attach request.

Qid is the servers representation of the file trees root.

See also: http://man.cat-v.org/plan_9/5/attach
*/
type RAttach struct {
	Qid qid.Qid
}

func (m *RAttach) encode(w io.Writer) (int64, error) {
	return m.Qid.Encode(w)
}

func (m *RAttach) decode(r io.Reader) (int64, error) {
	return m.Qid.Decode(r)
}
