package localmemory

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/registry"
)

//Local 本地内存作为注册中心
var Local r.IRegistry = newLocalMemory()

type localMemory struct {
	closeCh  chan struct{}
	nodes    map[string]string
	seqValue int32
	lock     sync.RWMutex
}

func newLocalMemory() *localMemory {
	return &localMemory{
		closeCh: make(chan struct{}),
		nodes:   make(map[string]string),
	}
}

func (l *localMemory) Exists(path string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	np := r.Join(path)
	if _, ok := l.nodes[np]; ok {
		return true, nil
	}
	return false, nil
}

var once sync.Once

func (l *localMemory) print() {
	once.Do(func() {
		for k, v := range l.nodes {
			fmt.Printf("%-48s %-32s\n", k, v)
		}
	})
}

func (l *localMemory) GetValue(path string) (data []byte, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if v, ok := l.nodes[r.Join(path)]; ok {
		return []byte(v), 0, nil
	}

	return nil, 0, fmt.Errorf("节点[%s]不存在", path)

}
func (l *localMemory) Update(path string, data string, version int32) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.nodes[r.Join(path)]; ok {
		l.nodes[r.Join(path)] = data
		return nil
	}
	return fmt.Errorf("节点[%s]不存在", path)
}
func (l *localMemory) GetChildren(path string) (paths []string, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	paths = make([]string, 0, 1)
	npath := r.Join(path)
	exists := make(map[string]string)
	for k := range l.nodes {
		if strings.HasPrefix(k, npath) && len(k) > len(npath) {
			list := strings.Split(strings.Trim(k[len(npath):], "/"), "/")
			name := list[0]
			if _, ok := exists[name]; !ok {
				exists[name] = name
				paths = append(paths, name)
			}

		}
	}
	return paths, 0, nil
}

func (l *localMemory) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	v := &eventWatcher{
		watcher: make(chan registry.ValueWatcher),
	}

	return v.watcher, nil

}
func (l *localMemory) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	return nil, nil
}
func (l *localMemory) Delete(path string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	np := r.Join(path)
	for k := range l.nodes {
		if strings.HasPrefix(k, np) {
			delete(l.nodes, k)
		}
	}
	return nil
}

func (l *localMemory) CreatePersistentNode(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	list := strings.Split(r.Join(path), "/")
	for i := range list {
		path := r.Join(list[:i]...)
		if _, ok := l.nodes[path]; !ok {
			l.nodes[path] = "{}"
		}
	}
	l.nodes[r.Join(path)] = data
	return nil
}
func (l *localMemory) CreateTempNode(path string, data string) (err error) {
	return l.CreatePersistentNode(path, data)
}
func (l *localMemory) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&l.seqValue, 1)
	rpath = fmt.Sprintf("%s%d", path, nid)
	err = l.CreatePersistentNode(rpath, data)
	return
}

func (l *localMemory) Close() error {
	return nil
}

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}
type valuesEntity struct {
	values  []string
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

func (v *valuesEntity) GetValue() ([]string, int32) {
	return v.values, v.version
}
func (v *valuesEntity) GetError() error {
	return v.Err
}
func (v *valuesEntity) GetPath() string {
	return v.path
}

type eventWatcher struct {
	watcher chan registry.ValueWatcher
}

//zkRegistry 基于zookeeper的注册中心
type lmFactory struct{}

//Build 根据配置生成文件系统注册中心
func (z *lmFactory) Create(addrs []string, u string, p string, log logger.ILogging) (r.IRegistry, error) {
	return Local, nil
}

func init() {
	r.Register(r.LocalMemory, &lmFactory{})
}
