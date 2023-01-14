package message // import "go.rbn.im/neinp/message"

import (
	"go.rbn.im/neinp/basic"
	"io"
)

/*9p defined error strings.
See: https://github.com/0intro/plan9/blob/7524062cfa4689019a4ed6fc22500ec209522ef0/sys/src/lib9p/srv.c#L10 */
const (
	BadAttachErrorString    = "unknown specifier in attach"
	BadOffsetErrorString    = "bad offset"
	BadCountErrorString     = "bad count"
	BotchErrorString        = "9P protocol botch"
	CreateNonDirErrorString = "create in non-directory"
	DupFidErrorString       = "duplicate fid"
	DupTagErrorString       = "duplicate tag"
	IsDirErrorString        = "is a directory"
	NoCreateErrorString     = "create prohibited"
	NoMemErrorString        = "out of memory"
	NoRemoveErrorString     = "remove prohibited"
	NoStatErrorString       = "stat prohibited"
	NotFoundErrorString     = "file not found"
	NoWriteErrorString      = "write prohibited"
	NoWstatErrorString      = "wstat prohibited"
	PermErrorString         = "permission denied"
	UnknownFidErrorString   = "unknown fid"
	BadDirErrorString       = "bad directory in wstat"
	WalkNoDirErrorString    = "walk in non-directory"
)

/*RError is the response to failed requests.

It usually contains one of the standard *ErrorString constants.

See also: http://man.cat-v.org/plan_9/5/error */
type RError struct {
	Ename string
}

func (m *RError) encode(w io.Writer) (int64, error) {
	return basic.StringEncode(w, m.Ename)
}

func (m *RError) decode(r io.Reader) (int64, error) {
	ename, n, err := basic.StringDecode(r)
	if err != nil {
		return n, err
	}

	m.Ename = ename
	return n, err
}
