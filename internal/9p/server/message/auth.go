package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/qid"
	"io"
)

/*TAuth facilitates authentication.

Afid is a new fid to be used for authentication.
Uname and Aname must be the same as for the following TAttach.

Authentication is performed by reading and writing the file referenced
by Afid, after receiving RAuth.
When complete Afid is used for authentication in TAttach.

The same Afid, Uname and Aname may be used for multiple TAttach messages.

If no authentication is required, RError is returned instead of
RAuth.

See also: http://man.cat-v.org/plan_9/5/attach
*/
type TAuth struct {
	Afid  fid.Fid
	Uname string
	Aname string
}

func (m *TAuth) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Afid))
	if err != nil {
		return n1, err
	}

	n2, err := basic.StringEncode(w, m.Uname)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := basic.StringEncode(w, m.Aname)
	if err != nil {
		return n1 + n2 + n3, err
	}

	return n1 + n2 + n3, nil
}

func (m *TAuth) decode(r io.Reader) (int64, error) {
	afid, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	uname, n2, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2, err
	}

	aname, n3, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	m.Afid = fid.Fid(afid)
	m.Uname = uname
	m.Aname = aname

	return n1 + n2 + n3, nil
}

/*RAuth is the response to a TAuth

Aqid is a file of type QidTypeAuth that can be read and
written using read and write messages to perform authentication.
The authentication protocol is not part of 9p.

See also: http://man.cat-v.org/plan_9/5/attach
*/
type RAuth struct {
	Aqid qid.Qid
}

func (m *RAuth) encode(w io.Writer) (int64, error) {
	return m.Aqid.Encode(w)
}

func (m *RAuth) decode(r io.Reader) (int64, error) {
	return m.Aqid.Decode(r)
}
