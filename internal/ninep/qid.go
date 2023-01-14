package ninep

import "fmt"

// bits in QID.Kind, from sys/include/libc.h
const (
	QTDIR    = 0x80 // type bit for directories
	QTAPPEND = 0x40 // type bit for append only files
	QTEXCL   = 0x20 // type bit for exclusive use files
	QTMOUNT  = 0x10 // type bit for mounted channel
	QTAUTH   = 0x08 // type bit for authentication file
	QTTMP    = 0x04 // type bit for not-backed-up file
	QTFILE   = 0x00 // plain file
)

// QID in Plan9 is defined in libc.h
type QID struct {
	Kind uint8  // uchar (normally called "Type")
	Vers uint32 // ulong
	Path uint64 // uvlong
}

func (q QID) String() string { return fmt.Sprintf("{0x%016x %d %d}", q.Path, q.Vers, q.Kind) }

func (q QID) IsDirectory() bool {
	return (q.Kind & QTDIR) != 0
}
