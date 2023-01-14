// Package id provides utility functions to convert unix uid/gid to strings.
package id // import "go.rbn.im/neinp/id"

import (
	"os/user"
	"strconv"
)

//UidToName returns the user name for an uid.
func UidToName(uid uint32) (string, error) {
	x, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return "", err
	}
	return x.Name, nil
}

//GidToName returns the group name for an gid.
func GidToName(gid uint32) (string, error) {
	y, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return "", err
	}
	return y.Name, nil
}
