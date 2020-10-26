package mocks

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
)

//Local 本地内存作为注册中心
var _ r.IRegistry = &TestRegistry{}

type TestRegistry struct {
	closeCh  chan struct{}
	nodes    map[string]string
	seqValue int32
	path     string
	lock     sync.RWMutex
	once     int
	Deep     int
}

//NewTestRegistry 构建基于文件系统的注册中心
func NewTestRegistry(platName string, systemName string, clusterName string, path string) *TestRegistry {
	f := &TestRegistry{
		closeCh: make(chan struct{}),
		nodes:   make(map[string]string),
		path:    path,
	}
	vnodes := map[string]map[string]interface{}{
		"api": map[string]interface{}{
			"server1": "value1",
			"server2": "value2",
			"server3": "value3",
			"server4": "value4",
			"server6": "value6",
			"main":    "main",
		},
	}

	for k, sub := range vnodes {
		for name, value := range sub {
			var path = r.Join(platName, systemName, k, clusterName, "hosts", name)
			if name == "main" {
				path = r.Join(platName, systemName, k, clusterName, "hosts")
			}
			buff, _ := json.Marshal(&value)
			f.nodes[path] = string(buff)
		}
	}
	f.nodes["/platname/apiserver/api/test/hosts1"] = "getErr"
	return f
}

func (l *TestRegistry) Exists(path string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for k, _ := range l.nodes {
		if strings.HasPrefix(k, path) {
			return true, nil
		}
	}
	// if _, ok := l.nodes[r.Join(path)]; ok {
	// 	return true, nil
	// }
	return false, nil
}
func (l *TestRegistry) GetValue(path string) (data []byte, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if v, ok := l.nodes[r.Join(path)]; ok {
		if v == "getErr" {
			return data, 0, fmt.Errorf("获取不到节点的值")
		}
		return []byte(v), 0, nil
	}
	return nil, 0, fmt.Errorf("节点[%s]不存在", path)

}
func (l *TestRegistry) Update(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.nodes[r.Join(path)]; ok {
		l.nodes[r.Join(path)] = data
		return nil
	}
	return fmt.Errorf("节点[%s]不存在", path)
}
func (l *TestRegistry) GetChildren(path string) (paths []string, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	paths = make([]string, 0, 1)
	npath := r.Join(path)

	for k := range l.nodes {
		if strings.HasPrefix(k, npath) {
			lk := k[len(npath):]
			if strings.Count(lk, "/") == 1 {
				paths = append(paths, strings.Trim(lk, "/"))
			}
		}
	}
	return paths, 0, nil
}

func (l *TestRegistry) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	v := &eventWatcher{
		watcher: make(chan registry.ValueWatcher, 1),
	}
	if l.once == 0 {
		l.once++
		fmt.Println(r.Join(path))
		l.nodes[r.Join(path)] = "value1-1"
		v.watcher <- &valueEntity{path: r.Join(path), Err: nil, Value: []byte("value1-1"), version: 0}
	}
	return v.watcher, nil
}
func (l *TestRegistry) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {

	data = make(chan registry.ChildrenWatcher, 1)

	if (path == "/platname/apiserver/api/test/hosts/server6" && l.Deep == 1) ||
		(path == "/platname/apiserver/api/test/hosts" && l.Deep == 2) ||
		(path == "/platname/apiserver/api/test" && l.Deep == 3) {
		if l.once == 0 {
			l.once++
			c, _, _ := l.GetChildren(path)
			npath := r.Join(path, "addnode")
			l.nodes[npath] = "addvalue"

			c = append(c, "addnode")
			data <- &valuesEntity{path: path, Err: nil, values: c, version: 0}
			return
		}
	}

	return data, nil
}
func (l *TestRegistry) Delete(path string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.nodes, r.Join(path))
	return nil
}

func (l *TestRegistry) CreatePersistentNode(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.nodes[r.Join(path)] = data
	return nil
}
func (l *TestRegistry) CreateTempNode(path string, data string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.nodes[r.Join(path)] = data
	return nil
}
func (l *TestRegistry) CreateSeqNode(path string, data string) (rpath string, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	nid := atomic.AddInt32(&l.seqValue, 1)
	rpath = fmt.Sprintf("%s%d", path, nid)
	l.nodes[rpath] = data
	return rpath, nil
}

func (l *TestRegistry) Close() error {
	return nil
}

type eventWatcher struct {
	watcher chan registry.ValueWatcher
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
