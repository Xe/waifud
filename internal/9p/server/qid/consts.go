package qid // import "go.rbn.im/neinp/qid"

//Type represents the type this qid symbolizes.
type Type uint8

const (
	// TypeFile is a plain file
	TypeFile Type = 0x00

	// TypeSymlink is a symbolic link
	TypeSymlink Type = 0x02

	// TypeTmp is a non-backed-up file
	TypeTmp Type = 0x04

	// TypeAuth is an authentication file
	TypeAuth Type = 0x08

	// TypeMount is a mounted channel
	TypeMount Type = 0x10

	// TypeExcl is a exclusive use file
	TypeExcl Type = 0x20

	// TypeAppend is a append only file
	TypeAppend Type = 0x40

	// TypeDir is a directory
	TypeDir Type = 0x80
)
