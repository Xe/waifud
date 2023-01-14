package client

import "flag"

// Flags
var (
	debugLog = flag.Bool("9p.debug", false, "enable debug logging of requests and responses")
)
