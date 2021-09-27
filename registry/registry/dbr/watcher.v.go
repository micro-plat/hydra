package dbr

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/registry"
)

type valueWatchers struct {
	db         dbs.IDB
	lk         sync.Mutex
	sqltexture *sqltexture
	watchers   map[string][]chan registry.ValueWatcher
	pths       cmap.ConcurrentMap
	closeCh    chan struct{}
	once       sync.Once
}

func newValueWatchers(db dbs.IDB, sqltexture *sqltexture) *valueWatchers {
	return &valueWatchers{
		db:         db,
		sqltexture: sqltexture,
		watchers:   make(map[string][]chan registry.ValueWatcher),
		pths:       cmap.New(2),
		closeCh:    make(chan struct{}),
	}
}
func (v *valueWatchers) Start() {
	tk := time.Tick(time.Second * 2)
	for {
		select {
		case <-tk:
			path := v.pths.Keys()
			for _, p := range path {
				data, err := v.db.Query(v.sqltexture.getValueChange, newInputByWatch(3, p))
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
func (v *valueWatchers) Watch(path string) chan registry.ValueWatcher {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		v.watchers[path] = make([]chan registry.ValueWatcher, 0, 1)
	}
	v.pths.SetIfAbsent(path, 0)
	ch := make(chan registry.ValueWatcher, 1)
	v.watchers[path] = append(v.watchers[path], ch)
	return ch
}

//通知节点值变化
func (v *valueWatchers) Notify(path string, version int32, value string) {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		return
	}
	for _, ch := range v.watchers[path] {
		v.pths.Remove(path)
		ch <- &valueEntity{Value: []byte(value), version: version, path: path}
		close(ch)
	}
	v.watchers[path] = make([]chan registry.ValueWatcher, 0)
}

//通知节点错误
func (v *valueWatchers) Error(path string, err error) {
	v.lk.Lock()
	defer v.lk.Unlock()
	if _, ok := v.watchers[path]; !ok {
		return
	}
	for _, ch := range v.watchers[path] {
		v.pths.Remove(path)
		ch <- &valueEntity{Err: err}
		close(ch)
	}
	v.watchers[path] = make([]chan registry.ValueWatcher, 0)
}
func (v *valueWatchers) Close() {
	v.once.Do(func() {
		close(v.closeCh)
	})
}

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}

func (v *valueEntity) GetPath() string {
	return v.path
}
func (v *valueEntity) GetValue() ([]byte, int32) {
	return v.Value, v.version
}
func (v *valueEntity) GetError() error {
	return v.Err
}
