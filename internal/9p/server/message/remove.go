package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"go.rbn.im/neinp/fid"
	"io"
)

/*TRemove signals the server to remove the file and clunk the Fid.

The Fid is clunked even if the remove fails. The request will fail
if the client doesn't have the permission to write the parent directory.

See also: http://man.cat-v.org/inferno/5/remove
*/
type TRemove struct {
	Fid fid.Fid
}

func (m *TRemove) encode(w io.Writer) (int64, error) {
	n, err := basic.Uint32Encode(w, uint32(m.Fid))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *TRemove) decode(r io.Reader) (int64, error) {
	f, n, err := basic.Uint32Decode(r)
	if err != nil {
		return n, err
	}

	m.Fid = fid.Fid(f)
	return n, nil
}

/*RRemove is the answer for a successful TRemove request.

See also: http://man.cat-v.org/inferno/5/remove
*/
type RRemove struct{}

func (m *RRemove) encode(w io.Writer) (int64, error) {
	return 0, nil
}

func (m *RRemove) decode(r io.Reader) (int64, error) {
	return 0, nil
}
