/*Package fs provides helpers to implement 9p objects: directories and files.*/
package fs // import "go.rbn.im/neinp/fs"

import (
	"go.rbn.im/neinp/message"
	"go.rbn.im/neinp/qid"
	"go.rbn.im/neinp/stat"
	"errors"
	"io"
)

//Entry is the general interface for 9p objects.
type Entry interface {
	Parent() Entry              // Parent returns this objects parent, it may be nil for the root object.
	Qid() qid.Qid               // Qid returns this files qid.
	Stat() stat.Stat            // Stat returns this files stat.
	Open() error                // Open prepares this file for access.
	Walk(string) (Entry, error) // Walk to a child or itself.
	io.ReadSeeker               // All files and directories need this.
}

//Dir implements a directory.
type Dir struct {
	parent   Entry
	stat     stat.Stat
	children []Entry
	*stat.Reader
}

//NewDir creates a new Dir ready for use.
//
//The stat and child entries are expected to be prepared by the caller.
//Children will have their parent set to this Dir.
func NewDir(s stat.Stat, children []Entry) *Dir {
	d := &Dir{
		stat: s,
	}

	for _, child := range children {
		switch t := child.(type) {
		case *Dir:
			t.parent = d
		case *File:
			t.parent = d
		}
	}

	d.children = children

	// just initialize with an empty stat.Reader here so we dont' hit nil
	d.Reader = stat.NewReader()

	return d
}

//Parent returns the parent of this dir, which may be nil if it is the root.
func (d *Dir) Parent() Entry {
	return d.parent
}

//Qid returns the qid of this dir.
func (d *Dir) Qid() qid.Qid {
	return d.stat.Qid
}

//Stat returns the stat of this dir.
func (d *Dir) Stat() stat.Stat {
	return d.stat
}

//Open prepares the dir for access.
//
//Internally, a stat.Reader is prepared from the contents of the children.
func (d *Dir) Open() error {
	stats := make([]stat.Stat, len(d.children))
	for i, child := range d.children {
		stats[i] = child.Stat()
	}

	d.Reader = stat.NewReader(stats...)
	return nil
}

//Walk to a child, the parent, or itself.
//
//wname is either the name of a child, "..", or "".
//If it is the name of a child, the respective entry is returned.
//For ".." the parent, or the dir itself (if it is the root) is returned.
//When wname is the empty string, the dir itself is returned.
func (d *Dir) Walk(wname string) (Entry, error) {
	if wname == ".." {
		if d.parent == nil {
			return d, nil
		}
		return d.parent, nil
	}

	for _, child := range d.children {
		if child.Stat().Name == wname {
			return child, nil
		}
	}

	return nil, errors.New(message.NotFoundErrorString)
}

//File implements a file.
type File struct {
	parent Entry
	stat   stat.Stat

	// this is the minimal interface we need, as a 9p read can also seek
	io.ReadSeeker
}

//NewFile prepares a new File ready for use.
//
//The stat and ReadSeeker are expected to be prepared by the caller.
func NewFile(s stat.Stat, rs io.ReadSeeker) *File {
	f := &File{
		stat:       s,
		ReadSeeker: rs,
	}

	return f
}

//Parent returns the parent.
func (f *File) Parent() Entry {
	return f.parent
}

//Qid returns the qid.
func (f *File) Qid() qid.Qid {
	return f.stat.Qid
}

//Stat returns the stat.
func (f *File) Stat() stat.Stat {
	return f.stat
}

//Walk can only walk to itself for plain files.
//
//It is used for the side effect of creating another fid for this file.
func (f *File) Walk(wname string) (Entry, error) {
	if len(wname) != 0 {
		return nil, errors.New(message.WalkNoDirErrorString)
	}

	return f, nil
}

//Open the file for access.
//
//This is a nop here and to be overriden by embedders.
func (f *File) Open() error {
	return nil
}
