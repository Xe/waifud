package neinp

import (
	"context"
	"errors"
	iofs "io/fs"
	"strings"

	"github.com/Xe/waifud/internal/9p/server/fid"
	"github.com/Xe/waifud/internal/9p/server/fs"
	"github.com/Xe/waifud/internal/9p/server/message"
)

type IOFS struct {
	NopP2000
	root fs.Dir
	fsys iofs.FS
	fids *fid.Map
}

func (i *IOFS) Version(ctx context.Context, m message.TVersion) (message.RVersion, error) {
	if !strings.HasPrefix(m.Version, "9P2000") {
		return message.RVersion{}, errors.New(message.BotchErrorString)
	}

	return message.RVersion{Version: m.Version, Msize: m.Msize}, nil
}

var _ P2000 = &IOFS{}
