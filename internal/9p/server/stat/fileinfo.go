package stat // import "go.rbn.im/neinp/stat"

import (
	"os"

	"go.rbn.im/neinp/qid"
)

// FileInfo creates Stat using os.FileInfo.Sys(). If using the information
// returned by Sys() fails, it returns a stat like returned by GenericStat.
func FileInfo(fi os.FileInfo) Stat {
	return fileInfo(fi)
}

// Generic creates a neinp.Stat not using FileInfo.Sys().
func Generic(fi os.FileInfo) Stat {
	size := fi.Size()
	if fi.IsDir() {
		size = 0
	}

	uid := "nobody"
	gid := "nogroup"

	stat := Stat{
		Qid:    qid.FileInfo(fi),
		Mode:   NeinMode(fi.Mode()),
		Atime:  fi.ModTime(),
		Mtime:  fi.ModTime(),
		Length: uint64(size),
		Name:   fi.Name(),
		Uid:    uid,
		Gid:    gid,
		Muid:   uid,
	}

	return stat
}
