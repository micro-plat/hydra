package dbr

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type tmpNodeWatchers struct {
	pths    cmap.ConcurrentMap
	db      dbs.IDB
	lk      sync.Mutex
	closeCh chan struct{}
	once    sync.Once
}

func newTmpNodeWatchers(db dbs.IDB) *tmpNodeWatchers {
	return &tmpNodeWatchers{
		db:      db,
		pths:    cmap.New(2),
		closeCh: make(chan struct{}),
	}
}

func (v *tmpNodeWatchers) Append(path string) {
	v.pths.SetIfAbsent(path, 0)
}
func (v *tmpNodeWatchers) Start() {
	tk := time.Tick(time.Second * 3)
	for {
		select {
		case <-tk:
			path := v.pths.Keys()
			v.db.Execute(aclUpdate, newInputByWatch(3, path...))
		case <-v.closeCh:
			return
		}
	}
}
func (v *tmpNodeWatchers) Close() {
	v.once.Do(func() {
		close(v.closeCh)
	})
}
