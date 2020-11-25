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

var _ r.IRegistry = &fs{}

type eventWatcher struct {
	watcher chan registry.ValueWatcher
	event   chan fsnotify.Event
}

type fs struct {
	watcher      *fsnotify.Watcher
	watcherMaps  map[string]*eventWatcher
	watchLock    sync.Mutex
	tempNode     []string
	tempNodeLock sync.Mutex
	seqNode      int32
	closeCh      chan struct{}
	prefix       string
}

func NewFileSystem(prefix string) (*fs, error) {
	// if err := checkPrivileges(); err != nil {
	// 	return nil, err
	// }
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fs{
		prefix:      strings.TrimRight(prefix, "/"),
		watcher:     w,
		watcherMaps: make(map[string]*eventWatcher),
		tempNode:    make([]string, 0, 2),
		seqNode:     10000,
		closeCh:     make(chan struct{}),
	}, nil
}

//Start 启动文件监控
func (l *fs) Start() {
	go func() {
	LOOP:
		for {
			select {
			case <-l.closeCh:
				break LOOP
			case event := <-l.watcher.Events:
				func(event fsnotify.Event) {
					l.watchLock.Lock()
					path := l.formatPath(event.Name)
					watcher, ok := l.watcherMaps[path]
					l.watchLock.Unlock()
					if !ok {
						return
					}
					watcher.event <- event
					delete(l.watcherMaps, path)
				}(event)

			}
		}
		l.watcher.Close()
	}()
}
func (l *fs) formatPath(path string) string {
	if !strings.HasPrefix(path, l.prefix) {
		return l.prefix + r.Join("/", path)
	}
	return path
}

func (l *fs) getRealPath(path string) string {
	if strings.HasSuffix(path, ".init") {
		return path
	}
	return fmt.Sprintf("%s/.init", path)
}

func (l *fs) Exists(path string) (bool, error) {
	_, err := os.Stat(l.formatPath(path))
	return err == nil || os.IsExist(err), nil
}
func (r *fs) getPaths(path string) []string {
	nodes := strings.Split(strings.Trim(path, "/"), "/")
	len := len(nodes)
	paths := make([]string, 0, len)
	for i := 0; i < len; i++ {
		npath := "/" + strings.Join(nodes[:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}

func (l *fs) GetValue(path string) (data []byte, version int32, err error) {
	rpath := l.formatPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		if strings.HasSuffix(rpath, ".init") {
			return []byte{}, 0, nil
		}
		return nil, 0, errors.New(rpath + "不存在")
	}
	if !fs.IsDir() {
		data, err = ioutil.ReadFile(rpath)
		version = int32(fs.ModTime().Unix())
		return
	}
	return l.GetValue(r.Join(path, ".init"))
}

func (l *fs) Update(path string, data string) (err error) {
	if b, _ := l.Exists(path); !b {
		return errors.New(path + "不存在")
	}

	rpath := l.formatPath(path)
	return ioutil.WriteFile(l.getRealPath(rpath), []byte(data), 0666)

}
func (l *fs) GetChildren(path string) (paths []string, version int32, err error) {
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
		if strings.HasSuffix(f.Name(), ".swp") || strings.HasPrefix(f.Name(), "~") || strings.HasPrefix(f.Name(), ".init") {
			continue
		}
		paths = append(paths, f.Name())
	}
	return paths, version, nil
}

func (l *fs) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	rpath := l.formatPath(path)
	absPath := rpath
	fs, _ := os.Stat(rpath)
	if fs != nil && fs.IsDir() {
		absPath = l.getRealPath(rpath)
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
				// fmt.Println("111111111111:", rpath)
				buff, version, err := l.GetValue(rpath)
				v.watcher <- &valueEntity{Value: buff, version: version, path: rpath, Err: err}
			default:
				// fmt.Println("22222222222:", rpath)
				v.watcher <- &valueEntity{path: rpath, Err: fmt.Errorf("文件发生变化:%v", event.Op)}
			}
		}
	}(rpath, l.watcherMaps[absPath])
	return l.watcherMaps[absPath].watcher, nil
}
func (l *fs) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	return nil, nil
}
func (l *fs) Delete(path string) error {

	if b, _ := l.Exists(path); !b {
		return nil
	}
	return os.RemoveAll(l.formatPath(path))
}

func (l *fs) CreatePersistentNode(path string, data string) (err error) {
	xpath := l.formatPath(path)
	paths := l.getPaths(path)
	for _, spath := range paths {
		rpath := l.formatPath(spath)
		rpath = l.getRealPath(rpath)
		if strings.HasPrefix(rpath, xpath) {
			if err = l.createNode(rpath, data); err != nil {
				return err
			}
		} else {
			if err = l.createNode(rpath, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *fs) createNode(path, data string) error {
	b, err := l.Exists(path)
	if err != nil {
		return err
	}
	if b {
		if data == "" {
			//节点有包含关系的时候  保证有数据的节点不被覆盖
			return nil
		}
		os.Remove(path)
	}
	if err = os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	f, err := os.Create(path) //创建文件
	if err != nil {
		return err
	}
	defer f.Close()
	err = os.Chmod(path, 0777)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(data); err != nil {
		return err
	}
	return nil
}
func (l *fs) CreateTempNode(path string, data string) (err error) {
	if err = l.CreatePersistentNode(path, data); err != nil {
		return err
	}
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	l.tempNode = append(l.tempNode, l.formatPath(path))
	return nil
}
func (l *fs) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&l.seqNode, 1)
	rpath = fmt.Sprintf("%s_%d", path, nid)
	return rpath, l.CreateTempNode(rpath, data)
}

func (l *fs) GetSeparator() string {
	return string(filepath.Separator)
}

func (l *fs) CanWirteDataInDir() bool {
	return false
}
func (l *fs) Close() error {
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	close(l.closeCh)
	for _, p := range l.tempNode {
		rp := l.getRealPath(p)
		if ok, _ := l.Exists(rp); ok {
			os.Remove(rp)
		}
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
