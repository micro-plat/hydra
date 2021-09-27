/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-27 16:01:56
 */
package dbr

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type tmpNodeWatchers struct {
	pths       cmap.ConcurrentMap
	db         dbs.IDB
	sqltexture *sqltexture
	lk         sync.Mutex
	closeCh    chan struct{}
	once       sync.Once
}

func newTmpNodeWatchers(db dbs.IDB, sqltexture *sqltexture) *tmpNodeWatchers {
	return &tmpNodeWatchers{
		db:         db,
		sqltexture: sqltexture,
		pths:       cmap.New(2),
		closeCh:    make(chan struct{}),
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
			v.db.Execute(v.sqltexture.aclUpdate, newInputBySelectIn(3, path...))
		case <-v.closeCh:
			return
		}
	}
}

func (v *tmpNodeWatchers) Close() {
	v.once.Do(func() {
		//退出时 删除所有的临时节点
		path := v.pths.Keys()
		v.db.Execute(v.sqltexture.clearTmpNode, newInputBySelectIn(0, path...))
		close(v.closeCh)
	})
}
