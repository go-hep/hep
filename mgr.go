package rio

import (
	"sync"
)

var (
	StreamMgr = streamMgr{
		db: make(map[string]*Stream),
	}

	RecordMgr = recordMgr{
		db: make(map[string]*Record),
	}

	BlockMgr = blockMgr{
		db: make(map[string]Block),
	}
)

type streamMgr struct {
	mu sync.RWMutex
	db map[string]*Stream
}

type recordMgr struct {
	mu sync.RWMutex
	db map[string]*Record
}

func (mgr *recordMgr) Get(name string) *Record {
	mgr.mu.RLock()
	rec, _ := mgr.db[name]
	mgr.mu.RUnlock()
	return rec
}

func (mgr *recordMgr) Has(name string) bool {
	mgr.mu.RLock()
	_, ok := mgr.db[name]
	mgr.mu.RUnlock()
	return ok
}

func (mgr *recordMgr) Add(name string) {
	mgr.mu.Lock()
	mgr.db[name] = &Record{name: name}
	mgr.mu.Unlock()
}

type blockMgr struct {
	mu sync.RWMutex
	db map[string]Block
}

// EOF
