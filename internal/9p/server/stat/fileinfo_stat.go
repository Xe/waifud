//+build !darwin,!plan9

package stat // import "go.rbn.im/neinp/stat"

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"go.rbn.im/neinp/qid"
)

// fileInfo creates Stat using os.FileInfo.Sys(). If using the information
// returned by Sys() fails, it returns a stat like returned by GenericStat.
func fileInfo(fi os.FileInfo) Stat {
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return Generic(fi)
	}

	size := fi.Size()
	if fi.IsDir() {
		size = 0
	}

	var uid, gid string

	x, err := user.LookupId(strconv.Itoa(int(s.Uid)))
	if err != nil {
		return Generic(fi)
	}
	uid = x.Name

	y, err := user.LookupGroupId(strconv.Itoa(int(s.Gid)))
	if err != nil {
		return Generic(fi)
	}
	gid = y.Name

	stat := Stat{
		Qid:    qid.FileInfo(fi),
		Mode:   Mode(fi.Mode()),
		Atime:  time.Unix(s.Atim.Sec, s.Atim.Nsec),
		Mtime:  time.Unix(s.Mtim.Sec, s.Mtim.Nsec),
		Length: uint64(size),
		Name:   fi.Name(),
		Uid:    uid,
		Gid:    gid,
		Muid:   uid,
	}

	return stat
}
