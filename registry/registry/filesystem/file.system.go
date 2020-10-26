package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/BurntSushi/toml"
	"github.com/micro-plat/hydra/global"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
)

//Local 本地内存作为注册中心
var _ r.IRegistry = &fileSystem{}

type fileSystem struct {
	closeCh  chan struct{}
	nodes    map[string]string
	seqValue int32
	path     string
	lock     sync.RWMutex
}

//NewFileSystem 构建基于文件系统的注册中心
func NewFileSystem(platName string, systemName string, clusterName string, path string) (*fileSystem, error) {
	f := &fileSystem{
		closeCh: make(chan struct{}),
		nodes:   make(map[string]string),
		path:    path,
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在:%s %v", path, err)
	}
	vnodes := make(map[string]map[string]interface{})
	if _, err := toml.DecodeFile(path, &vnodes); err != nil {
		return nil, err
	}
	for k, sub := range vnodes {
		for name, value := range sub {
			var path = r.Join(platName, systemName, k, clusterName, "conf", name)
			if name == "main" {
				path = r.Join(platName, systemName, k, clusterName, "conf")
			}
			buff, err := json.Marshal(&value)
			if err != nil {
				return nil, fmt.Errorf("转换配置信息为json串失败:%s %w", path, err)
			}
			f.nodes[path] = string(buff)
		}
	}
	return f, nil

}

func (l *fileSystem) Exists(path string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.nodes[r.Join(path)]; ok {
		return true, nil
	}
	return false, nil
}
func (l *fileSystem) GetValue(path string) (data []byte, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if v, ok := l.nodes[r.Join(path)]; ok {
		return []byte(v), 0, nil
	}
	return nil, 0, fmt.Errorf("节点[%s]不存在", path)

}
func (l *fileSystem) Update(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.nodes[r.Join(path)]; ok {
		l.nodes[r.Join(path)] = data
		return nil
	}
	return fmt.Errorf("节点[%s]不存在", path)
}
func (l *fileSystem) GetChildren(path string) (paths []string, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	paths = make([]string, 0, 1)
	npath := r.Join(path)
	for k := range l.nodes {
		lk := strings.TrimPrefix(k, npath)
		if len(lk) > 2 {
			paths = append(paths, strings.Trim(lk, "/"))
		}
	}
	return paths, 0, nil
}

func (l *fileSystem) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	v := &eventWatcher{
		watcher: make(chan registry.ValueWatcher),
	}

	return v.watcher, nil

}
func (l *fileSystem) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	return nil, nil
}
func (l *fileSystem) Delete(path string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.nodes, r.Join(path))
	return nil
}

func (l *fileSystem) CreatePersistentNode(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.nodes[r.Join(path)] = data
	return nil
}
func (l *fileSystem) CreateTempNode(path string, data string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.nodes[r.Join(path)] = data
	return nil
}
func (l *fileSystem) CreateSeqNode(path string, data string) (rpath string, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	nid := atomic.AddInt32(&l.seqValue, 1)
	rpath = fmt.Sprintf("%s%d", path, nid)
	l.nodes[rpath] = data
	return rpath, nil
}

func (l *fileSystem) Close() error {
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

//fsFactory 基于本地文件系统
type fsFactory struct {
	opts *r.Options
}

//Build 根据配置生成文件系统注册中心
func (z *fsFactory) Create(opts ...r.Option) (r.IRegistry, error) {
	//addrs []string, u string, p string, log logger.ILogging
	for i := range opts {
		opts[i](z.opts)
	}

	return NewFileSystem(global.Def.PlatName,
		global.Def.SysName,
		global.Def.ClusterName,
		filepath.Join(z.opts.Addrs[0], global.Def.LocalConfName))

}

func init() {
	r.Register(r.FileSystem, &fsFactory{
		opts: &r.Options{},
	})
}
