package client

// Every 9P message has an associated type indicating its contents and
// wire format:
const (
	Tversion = 100
	Rversion = 101
	Tauth    = 102
	Rauth    = 103
	Tattach  = 104
	Rattach  = 105
	/* There is no request for errors. */
	Rerror  = 107
	Tflush  = 108
	Rflush  = 109
	Twalk   = 110
	Rwalk   = 111
	Topen   = 112
	Ropen   = 113
	Tcreate = 114
	Rcreate = 115
	Tread   = 116
	Rread   = 117
	Twrite  = 118
	Rwrite  = 119
	Tclunk  = 120
	Rclunk  = 121
	Tremove = 122
	Rremove = 123
	Tstat   = 124
	Rstat   = 125
	Twstat  = 126
	Rwstat  = 127
	// Plan9 from User Space extensions
	Topenfd = 98
	Ropenfd = 99
)
