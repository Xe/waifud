package fid // import "go.rbn.im/neinp/fid"

import (
	"sync"
)

type Map struct {
	fids map[Fid]interface{}
	sync.RWMutex
}

func New() *Map {
	return &Map{fids: make(map[Fid]interface{})}
}

func (f *Map) Get(fid Fid) interface{} {
	f.RLock()
	defer f.RUnlock()

	v, ok := f.fids[fid]
	if !ok {
		return nil
	}

	return v
}

func (f *Map) Set(fid Fid, v interface{}) {
	f.Lock()
	defer f.Unlock()
	f.fids[fid] = v
}

func (f *Map) Delete(fid Fid) {
	f.Lock()
	defer f.Unlock()
	delete(f.fids, fid)
}
