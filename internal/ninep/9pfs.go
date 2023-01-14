package ninep

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
)

type file struct {
	FID    uint32
	cc     *clientConn
	offset int64
	iounit uint32
	QID    QID
}

func (f *file) Read(p []byte) (n int, err error) {
	n, err = f.ReadAt(p, f.offset)
	f.offset += int64(n)
	return n, err
}

func (f *file) ReadAt(p []byte, off int64) (n int, err error) {
	// Truncate read to iounit size if necessary.
	if uint32(len(p)) > f.iounit {
		p = p[:f.iounit]
	}
	count, err := f.cc.Read(context.TODO(), f.FID, uint64(off), p)
	if err != nil {
		return 0, err
	}
	if count == 0 && len(p) > 0 {
		return 0, io.EOF
	}
	return int(count), nil
}

func (f *file) Stat() (info os.FileInfo, err error) {
	stat, err := f.cc.Stat(context.TODO(), f.FID)
	return &statFileInfo{s: stat}, err
}

func (f *file) ReadDir(n int) (entries []fs.DirEntry, err error) {
	if !f.QID.IsDirectory() {
		return nil, errors.New("not a directory")
	}
	br := bufio.NewReader(f)
	unlimited := n <= 0
	for i := 0; i < n || unlimited; i++ {
		var stat Stat
		if err := readStat(br, &stat); err != nil {
			if unlimited && err == io.EOF {
				err = nil
			}
			return entries, err
		}
		entries = append(entries, &statFileInfo{s: stat})
	}
	return entries, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	var absOffset int64

	switch whence {
	case io.SeekStart:
		absOffset = offset
	case io.SeekCurrent:
		absOffset = f.offset + offset
	case io.SeekEnd:
		stat, err := f.Stat()
		if err != nil {
			return f.offset, err
		}
		absOffset = stat.Size() + offset
	default:
		return f.offset, fs.ErrInvalid
	}

	if absOffset < 0 {
		return f.offset, fs.ErrInvalid
	}
	f.offset = absOffset
	return f.offset, nil
}

func (f *file) Close() error {
	return f.cc.Clunk(context.TODO(), f.FID)
}

// TODO: Double check that the mode bits match.
type statFileInfo struct{ s Stat }

func (fi *statFileInfo) Name() string               { return fi.s.Name }
func (fi *statFileInfo) Size() int64                { return int64(fi.s.Length) }
func (fi *statFileInfo) Mode() fs.FileMode          { return os.FileMode(fi.s.Mode) }
func (fi *statFileInfo) ModTime() time.Time         { return time.Unix(int64(fi.s.Mtime), 0) }
func (fi *statFileInfo) IsDir() bool                { return (fi.s.Mode & ModeDir) != 0 }
func (fi *statFileInfo) Sys() interface{}           { return fi.s }
func (fi *statFileInfo) Type() fs.FileMode          { return fi.Mode().Type() }
func (fi *statFileInfo) Info() (fs.FileInfo, error) { return fi, nil }

type FS struct {
	cc      *clientConn
	nextFID uint32
	rootFID uint32
}

// Open opens a file for reading.
func (f *FS) Open(name string) (fs.File, error) {
	// TODO: Verify name format.
	components := strings.Split(name, "/")
	if len(name) == 0 {
		components = nil
	}

	// TODO: Track used FIDs instead of just cycling.
	f.nextFID++
	_, err := f.cc.Walk(context.TODO(), f.rootFID, f.nextFID, components)
	if err != nil {
		return nil, fmt.Errorf("9p walk: %w", err)
	}

	qid, iounit, err := f.cc.Open(context.TODO(), f.nextFID, ORead)
	if err != nil {
		return nil, fmt.Errorf("9p open: %w", err)
	}
	// If iounit is 0, we need to fall back to connection message
	// size - 24.
	if iounit == 0 {
		iounit = f.cc.msize - 24
	}

	return &file{FID: f.nextFID, cc: f.cc, iounit: iounit, QID: qid}, nil
}

var _ fs.FS = (*FS)(nil)
