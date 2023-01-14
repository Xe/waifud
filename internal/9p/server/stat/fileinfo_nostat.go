//+build darwin plan9

package stat // import "go.rbn.im/neinp/stat"

import (
	"os"
)

// fileInfo creates Stat using os.FileInfo.Sys(). If using the information
// returned by Sys() fails, it returns a stat like returned by GenericStat.
func fileInfo(fi os.FileInfo) Stat {
	return Generic(fi)
}
