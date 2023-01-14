package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"io"
)

/*TFlush signals the server that it should abort processing a message.

Instances of Server directly handle these message type by running the cancel function
of the context for the message which should be aborted. It is here for usage in
client implementations.

Unlike the real world, flushing never fails.

See also: http://man.cat-v.org/inferno/5/flush
*/
type TFlush struct {
	Oldtag uint16
}

func (m *TFlush) encode(w io.Writer) (int64, error) {
	return basic.Uint16Encode(w, m.Oldtag)
}

func (m *TFlush) decode(r io.Reader) (int64, error) {
	oldtag, n, err := basic.Uint16Decode(r)
	if err != nil {
		return n, err
	}

	m.Oldtag = oldtag
	return n, err
}

/*RFlush is the answer to a TFlush and signals that the flush has finished.

See also: http://man.cat-v.org/plan_9/5/flush */
type RFlush struct {
}

func (m *RFlush) encode(w io.Writer) (int64, error) {
	return 0, nil
}

func (m *RFlush) decode(r io.Reader) (int64, error) {
	return 0, nil
}
