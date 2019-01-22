package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
)

var _ r.IRegistry = &local{}

type eventWatcher struct {
	watcher chan registry.ValueWatcher
	event   chan fsnotify.Event
}

type local struct {
	watcher      *fsnotify.Watcher
	watcherMaps  map[string]*eventWatcher
	watchLock    sync.Mutex
	tempNode     []string
	tempNodeLock sync.Mutex
	seqNode      int32
	closeCh      chan struct{}
	prefix       string
}

func newLocal(prefix string) (*local, error) {
	if err := checkPrivileges(); err != nil {
		return nil, err
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &local{
		prefix:      strings.TrimRight(prefix, "/"),
		watcher:     w,
		watcherMaps: make(map[string]*eventWatcher),
		tempNode:    make([]string, 0, 2),
		seqNode:     10000,
		closeCh:     make(chan struct{}),
	}, nil
}

//Start 启动文件监控
func (l *local) Start() {
	go func() {
	LOOP:
		for {
			select {
			case <-l.closeCh:
				break LOOP
			case event := <-l.watcher.Events:
				func(event fsnotify.Event) {
					//fmt.Println("event:", event.Name, event.Op, event.String())
					l.watchLock.Lock()
					watcher, ok := l.watcherMaps[event.Name]
					l.watchLock.Unlock()
					if !ok {
						return
					}
					watcher.event <- event
					delete(l.watcherMaps, event.Name)
				}(event)

			}
		}
		l.watcher.Close()
	}()
}
func (l *local) formatPath(path string) string {
	if !strings.HasPrefix(path, l.prefix) {
		return l.prefix + r.Join("/", path)
	}
	return path
}
func (l *local) Exists(path string) (bool, error) {
	_, err := os.Stat(l.formatPath(path))
	return err == nil || os.IsExist(err), nil
}
func (l *local) GetValue(path string) (data []byte, version int32, err error) {
	rpath := l.formatPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return nil, 0, errors.New(rpath + "不存在")
	}
	if !fs.IsDir() {
		data, err = ioutil.ReadFile(rpath)
		version = int32(fs.ModTime().Unix())
		return
	}
	return l.GetValue(r.Join(path, ".init"))
}
func (l *local) Update(path string, data string, version int32) (err error) {
	if b, _ := l.Exists(path); !b {
		return errors.New(path + "不存在")
	}
	return ioutil.WriteFile(l.formatPath(path), []byte(data), 0666)

}
func (l *local) GetChildren(path string) (paths []string, version int32, err error) {
	rpath := l.formatPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return nil, 0, errors.New(path + "不存在")
	}
	version = int32(fs.ModTime().Unix())
	rf, err := ioutil.ReadDir(rpath)
	if err != nil {
		return nil, 0, err
	}
	paths = make([]string, 0, len(rf))
	for _, f := range rf {
		if strings.HasSuffix(f.Name(), ".swp") || strings.HasPrefix(f.Name(), "~") {
			continue
		}
		paths = append(paths, f.Name())
	}
	return paths, version, nil
}

func (l *local) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	rpath := l.formatPath(path)
	absPath := rpath
	fs, _ := os.Stat(rpath)
	if fs != nil && fs.IsDir() {
		absPath = r.Join(rpath, ".init")
	}
	l.watchLock.Lock()
	defer l.watchLock.Unlock()
	v, ok := l.watcherMaps[absPath]
	if ok {
		return v.watcher, nil
	}
	l.watcherMaps[absPath] = &eventWatcher{
		event:   make(chan fsnotify.Event),
		watcher: make(chan registry.ValueWatcher),
	}

	go func(rpath string, v *eventWatcher) {
		if err := l.watcher.Add(rpath); err != nil {
			v.watcher <- &valueEntity{path: rpath, Err: err}
		}
		select {
		case <-l.closeCh:
			return
		case event := <-v.event:
			switch event.Op {
			case fsnotify.Write, fsnotify.Create:
				buff, version, err := l.GetValue(rpath)
				v.watcher <- &valueEntity{Value: buff, version: version, path: rpath, Err: err}
			default:
				v.watcher <- &valueEntity{path: rpath, Err: fmt.Errorf("文件发生变化:%v", event.Op)}
			}
		}
	}(rpath, l.watcherMaps[absPath])
	return l.watcherMaps[absPath].watcher, nil
}
func (l *local) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	return nil, nil
}
func (l *local) Delete(path string) error {

	if b, _ := l.Exists(path); !b {
		return nil
	}
	return os.Remove(l.formatPath(path))
}

func (l *local) CreatePersistentNode(path string, data string) (err error) {
	rpath := l.formatPath(path)
	_, err = os.Stat(rpath)
	if err == nil || os.IsExist(err) {
		return fmt.Errorf("%s已存在", rpath)
	}
	if err = os.MkdirAll(filepath.Dir(rpath), 0777); err != nil {
		return err
	}
	f, err := os.Create(rpath) //创建文件
	if err != nil {
		return err
	}
	err = os.Chmod(rpath, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(data); err != nil {
		return err
	}
	return nil
}
func (l *local) CreateTempNode(path string, data string) (err error) {
	if err = l.CreatePersistentNode(path, data); err != nil {
		return err
	}
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	l.tempNode = append(l.tempNode, l.formatPath(path))
	return nil
}
func (l *local) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&l.seqNode, 1)
	rpath = r.Join(l.formatPath(path), fmt.Sprint(nid))
	return rpath, l.CreateTempNode(rpath, data)
}
func (l *local) GetSeparator() string {
	return string(filepath.Separator)
}

func (l *local) CanWirteDataInDir() bool {
	return false
}
func (l *local) Close() error {
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	close(l.closeCh)
	for _, p := range l.tempNode {
		os.Remove(p)
	}
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

func checkPrivileges() error {
	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return nil
			}
			return ErrRootPrivileges
		}
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	return fmt.Errorf("%v %s", ErrUnsupportedSystem, runtime.GOOS)
}

var ErrUnsupportedSystem = errors.New("Unsupported system")
var ErrRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")
