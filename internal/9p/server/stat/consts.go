package stat // import "go.rbn.im/neinp/stat"

import (
	"os"
)

/*Mode represents modes of a directory entry.*/
type Mode uint32

const (
	// SetGid is setgid (Unix, 9P2000.u)
	SetGid Mode = 0x00040000

	// SetUid is setuid (Unix, 9P2000.u)
	SetUid Mode = 0x00080000

	// Socket is a socket (Unix, 9P2000.u)
	Socket Mode = 0x00100000

	// NamedPipe is a named pipe (Unix, 9P2000.u)
	NamedPipe Mode = 0x00200000

	// Device is a device file (Unix, 9P2000.u)
	Device Mode = 0x00800000

	// Symlink is a symbolic link (Unix, 9P2000.u)
	Symlink Mode = 0x02000000

	// Tmp is a non-backed-up file
	Tmp Mode = 0x04000000

	// Auth is a authentication file
	Auth Mode = 0x08000000

	// Mount is a mounted channel
	Mount Mode = 0x10000000

	// Excl is a exclusive use file
	Excl Mode = 0x20000000

	// Append is a append only file
	Append Mode = 0x40000000

	// Dir is a directory
	Dir Mode = 0x80000000
)

//NeinMode converts os.FileMode to Mode.
func NeinMode(osmode os.FileMode) Mode {
	var mode Mode

	if osmode&os.ModeDir == os.ModeDir {
		mode |= Dir
	}

	if osmode&os.ModeAppend == os.ModeAppend {
		mode |= Append
	}

	if osmode&os.ModeExclusive == os.ModeExclusive {
		mode |= Excl
	}

	if osmode%os.ModeTemporary == os.ModeTemporary {
		mode |= Tmp
	}

	if osmode&os.ModeSymlink == os.ModeSymlink {
		mode |= Symlink
	}

	if osmode&os.ModeDevice == os.ModeDevice {
		mode |= Device
	}

	if osmode&os.ModeNamedPipe == os.ModeNamedPipe {
		mode |= NamedPipe
	}

	if osmode&os.ModeSocket == os.ModeSocket {
		mode |= Socket
	}

	if osmode&os.ModeSetgid == os.ModeSetgid {
		mode |= SetGid
	}

	if osmode&os.ModeSetuid == os.ModeSetuid {
		mode |= SetUid
	}

	mode |= Mode(uint32(osmode) & 0777)

	return mode
}

//OsMode translates Mode to os.FileMode.
func OsMode(mode Mode) os.FileMode {
	var osmode os.FileMode
	if mode&Dir == Dir {
		osmode |= os.ModeDir
	}

	if mode&Append == Append {
		osmode |= os.ModeAppend
	}

	if mode&Excl == Excl {
		osmode |= os.ModeExclusive
	}

	if mode&Tmp == Tmp {
		osmode |= os.ModeTemporary
	}

	if mode&Symlink == Symlink {
		osmode |= os.ModeSymlink
	}

	if mode&Device == Device {
		osmode |= os.ModeDevice
	}

	if mode&NamedPipe == NamedPipe {
		osmode |= os.ModeNamedPipe
	}

	if mode&Socket == Socket {
		osmode |= os.ModeSocket
	}

	if mode&SetGid == SetGid {
		osmode |= os.ModeSetgid
	}

	if mode&SetUid == SetUid {
		osmode |= os.ModeSetuid
	}

	osmode |= os.FileMode(mode & 0777)

	return os.FileMode(osmode)
}
