package rio

import "sync"

type streamManager struct {
	sync.RWMutex
	db map[string]*Stream
}

type recordManager struct {
	sync.RWMutex
	db map[string]*Record
}

type blockManager struct {
	sync.RWMutex
	db map[string]Block
}

var (
	StreamMgr = streamManager{
		db: make(map[string]*Stream),
	}

	RecordMgr = recordManager{
		db: make(map[string]*Record),
	}

	BlockMgr = blockManager{
		db: make(map[string]Block),
	}
)

// EOF
