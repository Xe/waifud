//Package qid provides the qid type and supporting functionality.
package qid // import "go.rbn.im/neinp/qid"

import (
	"go.rbn.im/neinp/basic"
	"io"
)

/*Qid represents the server's unique identification for a file.

Two files on the same server hierarchy are
the same if and only if their Qids are the same. (The
client may have multiple Fids pointing to a single file on a
server and hence having a single Qid.) The Qid fields hold a type,
specifying whether the file is a directory, append-only file, etc.,
the Qid Version and the Qid Path. The Path is an integer unique
among all files in the hierarchy.  If a file is deleted and recreated with
the same name in the same directory, the old and new path
components of the Qids should be different. The Version is
a version number for a file; typically, it is incremented
every time the file is modified.
*/
type Qid struct {
	Type    Type
	Version uint32
	Path    uint64
}

//Encode writes the 9p representation of the qid to an io.Writer.
func (q *Qid) Encode(w io.Writer) (int64, error) {
	n1, err := basic.Uint8Encode(w, uint8(q.Type))
	if err != nil {
		return n1, err
	}

	n2, err := basic.Uint32Encode(w, q.Version)
	if err != nil {
		return n1 + n2, err
	}

	n3, err := basic.Uint64Encode(w, q.Path)
	if err != nil {
		return n1 + n2 + n3, err
	}

	return n1 + n2 + n3, nil
}

//Decode reads the Qid data from an io.Reader.
func (q *Qid) Decode(r io.Reader) (int64, error) {
	typ, n1, err := basic.Uint8Decode(r)
	if err != nil {
		return n1, err
	}

	version, n2, err := basic.Uint32Decode(r)
	if err != nil {
		return n1 + n2, err
	}

	path, n3, err := basic.Uint64Decode(r)
	if err != nil {
		return n1 + n2 + n3, err
	}

	q.Type = Type(typ)
	q.Version = version
	q.Path = path

	return n1 + n2 + n3, err
}
