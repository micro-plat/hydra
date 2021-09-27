package dbr

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/registry"
)

type childrenWatchers struct {
	db         dbs.IDB
	lk         sync.Mutex
	sqltexture *sqltexture
	watchers   map[string][]chan registry.ChildrenWatcher
	pths       cmap.ConcurrentMap
	closeCh    chan struct{}
	once       sync.Once
}

func newChildrenWatchers(db dbs.IDB, sqltexture *sqltexture) *childrenWatchers {
	return &childrenWatchers{
		db:         db,
		sqltexture: sqltexture,
		watchers:   make(map[string][]chan registry.ChildrenWatcher),
		pths:       cmap.New(2),
		closeCh:    make(chan struct{}),
	}
}
func (v *childrenWatchers) Start() {
	tk := time.Tick(time.Second * 2)
	for {
		select {
		case <-tk:
			path := v.pths.Keys()
			fmt.Printf("xxxxxxx:%+v \n", path)
			for _, p := range path {
				data, err := v.db.Query(v.sqltexture.getChildrenChange, newInputByWatch(3, p))
				if err == nil {
					for _, r := range data {
						go v.Notify(r.GetString(FieldPath), r.GetInt32(FieldDataVersion), r.GetString(FieldValue))
					}
				}
			}

		case <-v.closeCh:
			return
		}
	}
}

//监控节点的值变化
func (v *childrenWatchers) Watch(path string) chan registry.ChildrenWatcher {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		v.watchers[path] = make([]chan registry.ChildrenWatcher, 0, 1)
	}
	v.pths.SetIfAbsent(path, 0)
	ch := make(chan registry.ChildrenWatcher, 1)
	v.watchers[path] = append(v.watchers[path], ch)
	return ch
}

//通知节点值变化
func (v *childrenWatchers) Notify(path string, version int32, paths ...string) {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		return
	}
	for _, ch := range v.watchers[path] {
		v.pths.Remove(path)
		ch <- &childrenEntity{version: version, path: path}
		close(ch)
	}
	v.watchers[path] = make([]chan registry.ChildrenWatcher, 0)
}

//通知节点错误
func (v *childrenWatchers) Error(path string, err error) {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		return
	}
	for _, ch := range v.watchers[path] {
		v.pths.Remove(path)
		ch <- &childrenEntity{Err: err}
		close(ch)
	}
	v.watchers[path] = make([]chan registry.ChildrenWatcher, 0)
}
func (v *childrenWatchers) Close() {
	v.once.Do(func() {
		close(v.closeCh)
	})
}

type childrenEntity struct {
	children []string
	version  int32
	path     string
	Err      error
}

func (v *childrenEntity) GetValue() ([]string, int32) {
	return v.children, v.version
}
func (v *childrenEntity) GetError() error {
	return v.Err
}
func (v *childrenEntity) GetPath() string {
	return v.path
}
