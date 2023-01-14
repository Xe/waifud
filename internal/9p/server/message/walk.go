package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"go.rbn.im/neinp/qid"
	"io"
)

/*TWalk asks to walk down a directory hierarchy.

Fid is an already existing fid, Newfid is a proposed new
fid that the client wants to use for the result of the walk.
Wname are successive path elements to walk to, and may be of
length zero. Fid must be a directory, only if Wname is of length
zero. In this case Newfid will represent the same file as Fid.

Fid must be valid and not already opened by a TOpen or TCreate
request.

Wname may be "..", which represents the parent directory,
"." is not used and forbidden, zero length Wname is used for this.

Wname should have a maximum length of 16 elements (this limit is
not enforced here).

See also: http://man.cat-v.org/inferno/5/walk
*/
type TWalk struct {
	Fid    fid.Fid
	Newfid fid.Fid
	Wname  []string
}

func (m *TWalk) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint32Encode(w, uint32(m.Newfid))
	if err != nil {
		return n1 + n2, err
	}

	nwname := uint16(len(m.Wname))
	n3, err := basic.Uint16Encode(w, nwname)
	if err != nil {
		return n1 + n2 + n3, err
	}

	var n4 int64
	for _, v := range m.Wname {
		n, err := basic.StringEncode(w, v)
		n4 += n
		if err != nil {
			return n1 + n2 + n3 + n4, err
		}
	}

	return n1 + n2 + n3 + n4, err
}

func (m *TWalk) decode(r io.Reader) (int64, error) {
	f, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	newf, n2, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	nwname, n3, err := basic.Uint16Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	wnames := make([]string, nwname)

	var n4 int64
	for i := uint16(0); i < nwname; i++ {
		wname, n, err := basic.StringDecode(r)
		n4 += n
		if err != nil {
			return n1 + n2 + n3 + n4, err
		}
		wnames[i] = wname
	}

	m.Fid = fid.Fid(f)
	m.Newfid = fid.Fid(newf)
	m.Wname = wnames

	return n1 + n2 + n3 + n4, nil
}

/*RWalk is the response to a TWalk.

Wqid contains the qids of the Wnames in a successfull walk.
If the walk fails at the first element an error is returned,
if it fails at a later point, the qids of the successfully walked
elements are returned (len(Wqid) < len(Wname).

See also: http://man.cat-v.org/inferno/5/walk
*/
type RWalk struct {
	Wqid []qid.Qid
}

func (m *RWalk) encode(w io.Writer) (int64, error) {
	nwqid := uint16(len(m.Wqid))
	n1, err := basic.Uint16Encode(w, nwqid)
	if err != nil {
		return n1, err
	}

	var n2 int64
	for _, v := range m.Wqid {
		n, err := v.Encode(w)
		n2 += n
		if err != nil {
			return n1 + n2, err
		}
	}

	return n1 + n2, nil
}

func (m *RWalk) decode(r io.Reader) (int64, error) {
	nwqid, n1, err := basic.Uint16Decode(r)
	if err != nil {
		return n1, err
	}

	wqids := make([]qid.Qid, nwqid)

	var n2 int64
	for i := uint16(0); i < nwqid; i++ {
		n, err := wqids[i].Decode(r)
		n2 += n
		if err != nil {
			return n1 + n2, err
		}
	}

	m.Wqid = wqids

	return n1 + n2, nil
}
