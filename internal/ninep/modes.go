package ninep

import (
	"strings"
)

// 9p Stat mode flags for file meta-information.
const (
	ModeDir    = 0x80000000
	ModeAppend = 0x40000000
	ModeExcl   = 0x20000000
	ModeMount  = 0x10000000
	ModeAuth   = 0x08000000
	ModeTmp    = 0x04000000
)

// 9P2000.u extensions to 9p Stat mode flags.
const (
	ModeUnixDev   = 0x00800000
	ModeSymlink   = 0x00400000
	ModeNamedPipe = 0x00200000
	ModeSocket    = 0x00100000
)

// The read, write and execute bits are stored in the three least
// significant octets of Stat.Mode, for user, group and others.
const (
	ModeUserRead   = 0400
	ModeUserWrite  = 0200
	ModeUserExec   = 0100
	ModeGroupRead  = 0040
	ModeGroupWrite = 0020
	ModeGroupExec  = 0010
	ModeOtherRead  = 0004
	ModeOtherWrite = 0002
	ModeOtherExec  = 0001
)

// ModeString builds a "ls" style mode string for the given mode value.
// For example, ModeString(0744) == "-rwxr--r--".
func ModeString(mode uint32) string {
	b := strings.Builder{}
	for _, s := range []struct {
		mask   uint32
		symbol byte
	}{
		{ModeDir, 'd'},
		{ModeUserRead, 'r'},
		{ModeUserWrite, 'w'},
		{ModeUserExec, 'x'},
		{ModeGroupRead, 'r'},
		{ModeGroupWrite, 'w'},
		{ModeGroupExec, 'x'},
		{ModeOtherRead, 'r'},
		{ModeOtherWrite, 'w'},
		{ModeOtherExec, 'x'},
	} {
		if mode&s.mask != 0 {
			b.WriteByte(s.symbol)
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}
