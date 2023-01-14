package message // import "go.rbn.im/neinp/message"

type messageType uint8

// Internal constants for tag values of messages.
const (
	tversion messageType = iota + 100
	rversion
	tauth
	rauth
	tattach
	rattach
	terror // unused, good name for that ;)
	rerror
	tflush
	rflush
	twalk
	rwalk
	topen
	ropen
	tcreate
	rcreate
	tread
	rread
	twrite
	rwrite
	tclunk
	rclunk
	tremove
	rremove
	tstat
	rstat
	twstat
	rwstat
)
