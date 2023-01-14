package qid // import "go.rbn.im/neinp/qid"

import (
	"hash/fnv"
	"os"
)

// FileInfo creates a neinp.Qid from os.FileInfo, using Sys of FileInfo if possible.
func FileInfo(fi os.FileInfo) Qid {
	return fileInfo(fi)
}

// Generic creates a neinp.Qid not using FileInfo.Sys().
func Generic(fi os.FileInfo) Qid {
	t := TypeFile
	if fi.IsDir() {
		t = TypeDir
	}

	sec := fi.ModTime().Unix()
	nsec := fi.ModTime().UnixNano()

	v := fnv.New32a()
	v.Write([]byte{uint8(sec), uint8(sec >> 8), uint8(sec >> 16), uint8(sec >> 24), uint8(sec >> 32), uint8(sec >> 40), uint8(sec >> 48), uint8(sec >> 56)})
	v.Write([]byte{uint8(nsec), uint8(nsec >> 8), uint8(nsec >> 16), uint8(nsec >> 24), uint8(nsec >> 32), uint8(nsec >> 40), uint8(nsec >> 48), uint8(nsec >> 56)})

	p := fnv.New64a()
	p.Write([]byte(fi.Name()))

	return Qid{Type: t, Version: v.Sum32(), Path: p.Sum64()}
}
