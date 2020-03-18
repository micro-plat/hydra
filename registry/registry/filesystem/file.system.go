package filesystem

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

var _ r.IRegistry = &fileSystem{}

type eventWatcher struct {
	watcher chan registry.ValueWatcher
	event   chan fsnotify.Event
}

type fileSystem struct {
	watcher      *fsnotify.Watcher
	watcherMaps  map[string]*eventWatcher
	watchLock    sync.Mutex
	tempNode     []string
	tempNodeLock sync.Mutex
	seqNode      int32
	closeCh      chan struct{}
	path         string
}

func newFileSystem(path string) (*fileSystem, error) {
	if err := checkPrivileges(); err != nil {
		return nil, err
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fileSystem{
		path:        path,
		watcher:     w,
		watcherMaps: make(map[string]*eventWatcher),
		tempNode:    make([]string, 0, 2),
		seqNode:     10000,
		closeCh:     make(chan struct{}),
	}, nil
}

//Start 启动文件监控
func (l *fileSystem) Start() {
	go func() {
	LOOP:
		for {
			select {
			case <-l.closeCh:
				break LOOP
			case event := <-l.watcher.Events:
				func(event fsnotify.Event) {
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
func (l *fileSystem) formatPath(path string) string {
	npath := strings.Replace(path, "/", string(os.PathSeparator), -1)
	if !strings.HasPrefix(path, l.path) {
		return filepath.Join(l.path, npath)
	}
	return npath
}
func (l *fileSystem) getStaticPath(path string) string {
	p := l.formatPath(path)
	if !strings.HasSuffix(path, ".init") {
		return filepath.Join(p, ".init")
	}
	return p
}
func (l *fileSystem) getSeqPath(path string) string {
	p := l.formatPath(path)
	if strings.HasSuffix(path, ".init") {
		return p
	}
	nid := atomic.AddInt32(&l.seqNode, 1)
	return l.getStaticPath(fmt.Sprintf("%s-%d", path, nid))
}
func (l *fileSystem) Exists(path string) (bool, error) {
	_, err := os.Stat(l.getStaticPath(path))
	return err == nil || os.IsExist(err), nil
}
func (l *fileSystem) GetValue(path string) (data []byte, version int32, err error) {
	rpath := l.getStaticPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return nil, 0, nil
	}
	data, err = ioutil.ReadFile(rpath)
	version = int32(fs.ModTime().Unix())
	return

}
func (l *fileSystem) Update(path string, data string, version int32) (err error) {
	rpath := l.getStaticPath(path)
	if b, _ := l.Exists(rpath); !b {
		return errors.New(rpath + "不存在")
	}
	return ioutil.WriteFile(rpath, []byte(data), 0666)

}
func (l *fileSystem) GetChildren(path string) (paths []string, version int32, err error) {
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

func (l *fileSystem) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	rpath := l.getStaticPath(path)
	l.watchLock.Lock()
	defer l.watchLock.Unlock()
	v, ok := l.watcherMaps[rpath]
	if ok {
		return v.watcher, nil
	}
	l.watcherMaps[rpath] = &eventWatcher{
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
	}(rpath, l.watcherMaps[rpath])
	return l.watcherMaps[rpath].watcher, nil
}
func (l *fileSystem) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	// rpath := l.formatPath(path)
	// fs, err := os.Stat(rpath)
	// if err != nil {
	// 	return nil, err
	// }

	// time := fs.ModTime()

	return nil, nil
}
func (l *fileSystem) Delete(path string) error {
	rpath := l.formatPath(path)
	if b, _ := l.Exists(rpath); !b {
		return nil
	}
	return os.Remove(rpath)
}

func (l *fileSystem) CreatePersistentNode(path string, data string) (err error) {
	rpath := l.getStaticPath(path)
	_, err = os.Stat(rpath)
	if err == nil || os.IsExist(err) {
		os.Remove(rpath)
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
func (l *fileSystem) CreateTempNode(path string, data string) (err error) {
	rpath := l.getStaticPath(path)
	if err = l.CreatePersistentNode(path, data); err != nil {
		return err
	}
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	l.tempNode = append(l.tempNode, rpath)
	return nil
}
func (l *fileSystem) CreateSeqNode(path string, data string) (rpath string, err error) {
	rpath = l.getSeqPath(path)
	return rpath, l.CreateTempNode(rpath, data)
}

func (l *fileSystem) Close() error {
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
			return errRootPrivileges
		}
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	return fmt.Errorf("%v %s", errUnsupportedSystem, runtime.GOOS)
}

var errUnsupportedSystem = errors.New("Unsupported system")
var errRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")
