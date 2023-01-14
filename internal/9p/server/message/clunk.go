package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"io"
)

/*TClunk requests to forget a fid.

The mapping on the server from fid to actual file should be removed,
not touching the actual file unless it was used OpenClose.

Regardless of the result of this call, the fid will not be valid anymore
and can be reused.

See also: http://man.cat-v.org/plan_9/5/clunk
*/
type TClunk struct {
	Fid fid.Fid
}

func (m *TClunk) encode(w io.Writer) (int64, error) {
	n, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *TClunk) decode(r io.Reader) (int64, error) {
	f, n, err := basic.Uint32Decode(r)
	if err != nil {
		return n, err
	}

	m.Fid = fid.Fid(f)

	return n, err
}

/*RClunk signals a successful clunk.

After the fid is successfully clunked it can be reused.

See also: http://man.cat-v.org/plan_9/5/clunk
*/
type RClunk struct{}

func (m *RClunk) encode(w io.Writer) (int64, error) {
	return 0, nil
}

func (m *RClunk) decode(r io.Reader) (int64, error) {
	return 0, nil
}
