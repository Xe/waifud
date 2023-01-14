//+build darwin plan9

package qid // import "go.rbn.im/neinp/qid"

import (
	"os"
)

// fileInfo creates a neinp.Qid from os.FileInfo on systems where syscall.Stat
// is missing or limited.
func fileInfo(fi os.FileInfo) Qid {
	return Generic(fi)
}
