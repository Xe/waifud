package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"io"
)

/*TVersion initializes a new connection.

Msize is the maximum message size the client will ever create, and
includes _all_ protocol data (also 9p type, tag and size fields).
As the contents of these fields are handled by neinp.Server, the msize
has to include the length of the size (4 octets) type (1 octet) and tag (2 octets) fields,
even if not visible to neinp.P2000 implementers. The msize returned with RVersion
will be read by neinp.Server to setup a io.LimitReader for message reception.

Version identifies the level of the protocol, it must always start with
"9P" (though not enforced here).

See also: http://man.cat-v.org/inferno/5/version
*/
type TVersion struct {
	Msize   uint32
	Version string
}

func (m *TVersion) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, m.Msize)
	if err != nil {
		return n1, err
	}

	n2, err := basic.StringEncode(w, m.Version)
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}

func (m *TVersion) decode(r io.Reader) (int64, error) {
	msize, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	version, n2, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2, err
	}

	m.Msize = msize
	m.Version = version
	return n1 + n2, nil
}

/*RVersion is the servers reply to a TVersionMessage.

Msize must be lower or equal of what the client requested.

Version may be set to the clients version string or a string
of an earlier protocol version. If the server doesn't understand
the requested version, it replies with the version string "unknown".

A version request starts a new session, so any remaining I/O operations
are aborted and allocated fids are clunked.

See also: http://man.cat-v.org/inferno/5/version
*/
type RVersion struct {
	Msize   uint32
	Version string
}

func (m *RVersion) encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint32Encode(w, m.Msize)
	if err != nil {
		return n1, err
	}

	n2, err := basic.StringEncode(w, m.Version)
	if err != nil {
		return n1 + n2, err
	}

	return n1 + n2, nil
}

func (m *RVersion) decode(r io.Reader) (int64, error) {
	msize, n1, err := basic.Uint32Decode(r)
	if err != nil {
		return n1, err
	}

	version, n2, err := basic.StringDecode(r)
	if err != nil {
		return n1 + n2, err
	}

	m.Msize = msize
	m.Version = version
	return n1 + n2, nil
}
