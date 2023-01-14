package neinp

import (
	"context"
	"errors"

	"github.com/Xe/waifud/internal/9p/server/message"
)

// NopP2000 is a dummy impementation of interface P2000, returning error for every call.
// It can be embedded into own implementations if some functions aren't required.
// For a working mount (with p9p 9pfuse or linux kernel module), at least Version, Attach,
// Stat, Walk, Open, Read, and Clunk need to be implemented.
type NopP2000 struct{}

var errNop = errors.New("not implemented")

func (f *NopP2000) Version(context.Context, message.TVersion) (message.RVersion, error) {
	return message.RVersion{}, errNop
}

func (f *NopP2000) Auth(context.Context, message.TAuth) (message.RAuth, error) {
	return message.RAuth{}, errNop
}

func (f *NopP2000) Attach(context.Context, message.TAttach) (message.RAttach, error) {
	return message.RAttach{}, errNop
}

func (f *NopP2000) Walk(context.Context, message.TWalk) (message.RWalk, error) {
	return message.RWalk{}, errNop
}

func (f *NopP2000) Open(context.Context, message.TOpen) (message.ROpen, error) {
	return message.ROpen{}, errNop
}

func (f *NopP2000) Create(context.Context, message.TCreate) (message.RCreate, error) {
	return message.RCreate{}, errNop
}

func (f *NopP2000) Read(context.Context, message.TRead) (message.RRead, error) {
	return message.RRead{}, errNop
}

func (f *NopP2000) Write(context.Context, message.TWrite) (message.RWrite, error) {
	return message.RWrite{}, errNop
}

func (f *NopP2000) Clunk(context.Context, message.TClunk) (message.RClunk, error) {
	return message.RClunk{}, errNop
}

func (f *NopP2000) Remove(context.Context, message.TRemove) (message.RRemove, error) {
	return message.RRemove{}, errNop
}

func (f *NopP2000) Stat(context.Context, message.TStat) (message.RStat, error) {
	return message.RStat{}, errNop
}

func (f *NopP2000) Wstat(context.Context, message.TWstat) (message.RWstat, error) {
	return message.RWstat{}, errNop
}

func (f *NopP2000) Close() error {
	return nil
}
