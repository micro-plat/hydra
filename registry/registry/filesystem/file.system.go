package filesystem

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	//"github.com/micro-plat/hydra/global/compatible"

	"github.com/fsnotify/fsnotify"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
)

var _ r.IRegistry = &fs{}
var fileMode = os.FileMode(0664)
var dirMode = os.FileMode(0755)

type eventWatcher struct {
	watcher  chan registry.ValueWatcher
	event    chan fsnotify.Event
	syncChan chan string
}

type fs struct {
	watcher      *fsnotify.Watcher
	watcherMaps  map[string]*eventWatcher
	watchLock    sync.Mutex
	tempNodes    map[string]bool
	tempNodeLock sync.Mutex
	closeCh      chan struct{}
	rootDir      string
	done         bool
}

//NewFileSystem 文件系统的注册中心
func NewFileSystem(rootDir string) (*fs, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fs{
		rootDir:     strings.TrimRight(rootDir, "/"),
		watcher:     w,
		watcherMaps: make(map[string]*eventWatcher),
		tempNodes:   make(map[string]bool),
		closeCh:     make(chan struct{}),
	}, nil
}

//Start 启动文件监控
func (l *fs) Start() {
	go func() {
		for {
			select {
			case <-l.closeCh:
				return
			case event := <-l.watcher.Events:
				if l.done {
					return
				}
				func(event fsnotify.Event) {
					l.watchLock.Lock()
					defer l.watchLock.Unlock()
					dataPath := l.formatPath(event.Name)
					path := filepath.Dir(dataPath)
					watcher, ok := l.watcherMaps[path]
					if !ok {
						return
					}
					watcher.event <- event
				}(event)

			}
		}
		l.watcher.Close()
	}()
}

//formatPath 将rootDir 构建到路径中去
func (l *fs) formatPath(path string) string {
	if !strings.HasPrefix(path, l.rootDir) {
		return l.rootDir + r.Join("/", path)
	}
	return path
}

//getDataPath 获取目录下的.init路径
func (l *fs) getDataPath(path string) string {
	if strings.HasSuffix(path, ".init") {
		return path
	}
	return fmt.Sprintf("%s/.init", path)
}

func (l *fs) Exists(path string) (bool, error) {
	_, err := os.Stat(l.formatPath(path))
	return err == nil || os.IsExist(err), nil
}

func (l *fs) GetValue(path string) (data []byte, version int32, err error) {

	rpath := l.formatPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return []byte{}, 0, nil
	}
	dataPath := l.getDataPath(rpath)
	fs, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		return []byte{}, 0, nil
	}
	data, err = ioutil.ReadFile(dataPath)
	version = int32(fs.ModTime().Unix())
	return

}

func (l *fs) Update(path string, data string) (err error) {
	if b, _ := l.Exists(path); !b {
		return errors.New(path + "不存在")
	}

	rpath := l.formatPath(path)
	return ioutil.WriteFile(l.getDataPath(rpath), []byte(data), fileMode)

}
func (l *fs) GetChildren(path string) (paths []string, version int32, err error) {
	rpath := l.formatPath(path)
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return nil, 0, errors.New(path + "不存在")
	}
	version = int32(fs.ModTime().Unix())
	children, err := ioutil.ReadDir(rpath)
	if err != nil {
		return nil, 0, err
	}
	paths = make([]string, 0, len(children))
	for _, f := range children {
		if strings.HasSuffix(f.Name(), ".swp") || strings.HasPrefix(f.Name(), "~") || strings.HasPrefix(f.Name(), ".init") {
			continue
		}
		paths = append(paths, f.Name())
	}
	return paths, version, nil
}

func (l *fs) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	realPath := l.formatPath(path)
	_, err = os.Stat(realPath)
	if os.IsNotExist(err) {
		err = fmt.Errorf("Watch path:%s 不存在", path)
		return
	}

	l.watchLock.Lock()
	defer l.watchLock.Unlock()
	v, ok := l.watcherMaps[realPath]
	if ok {
		return v.watcher, nil
	}
	l.watcherMaps[realPath] = &eventWatcher{
		event:    make(chan fsnotify.Event),
		watcher:  make(chan registry.ValueWatcher),
		syncChan: make(chan string, 100),
	}
	go func(rpath string, v *eventWatcher) {
		dataFile := l.getDataPath(rpath)
		if err := l.watcher.Add(dataFile); err != nil {
			v.watcher <- &valueEntity{path: rpath, Err: err}
		}
		go func(evtw *eventWatcher) {
			ticker := time.NewTicker(time.Second * 2)
			for {
				select {
				case <-ticker.C:
					path := ""
				INFOR:
					for {
						select {
						case p := <-evtw.syncChan:
							path = p
						default:
							break INFOR

						}
					}
					if len(path) > 0 {
						ett := &valueEntity{
							path: rpath,
						}
						evtw.watcher <- ett
					}
				}
			}
		}(v)

		for {
			select {
			case <-l.closeCh:
				return
			case event := <-v.event:
				if event.Op == fsnotify.Chmod || event.Op == fsnotify.Rename {
					break
				}
				v.syncChan <- rpath
			}
		}
	}(realPath, l.watcherMaps[realPath])

	return l.watcherMaps[realPath].watcher, nil
}

func (l *fs) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	return nil, nil
}

func (l *fs) Delete(path string) error {
	return os.RemoveAll(l.formatPath(path))
}

func (l *fs) CreatePersistentNode(path string, data string) (err error) {
	err = l.createDirPath(path)
	if err != nil {
		return
	}
	err = l.createNodeData(path, data)
	if err != nil {
		return
	}
	return nil
}

func (l *fs) createDirPath(path string) error {
	realPath := l.formatPath(path)
	_, err := os.Stat(realPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(realPath, dirMode)
	}
	return nil
}

func (l *fs) createNodeData(path, data string) error {
	realPath := l.formatPath(path)
	dataPath := l.getDataPath(realPath)
	return ioutil.WriteFile(dataPath, []byte(data), fileMode)
}

func (l *fs) CreateTempNode(path string, data string) (err error) {
	if err = l.CreatePersistentNode(path, data); err != nil {
		return err
	}
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	l.tempNodes[l.formatPath(path)] = true
	return nil
}
func (l *fs) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := time.Now().UnixNano()
	rpath = fmt.Sprintf("%s_%d", path, nid)
	return rpath, l.CreateTempNode(rpath, data)
}

func (l *fs) GetSeparator() string {
	return string(filepath.Separator)
}

func (l *fs) Close() error {
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	if l.done {
		return nil
	}
	l.done = true
	close(l.closeCh)
	for path := range l.tempNodes {
		os.RemoveAll(path)
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
