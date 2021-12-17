package nfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//local 本地文件管理
type local struct {
	path        string
	fpPath      string
	currentAddr string
	FPS         cmap.ConcurrentMap
	once        sync.Once
	fsWatcher   *fsnotify.Watcher
	init        bool
	nfsChecker  sync.WaitGroup
	lockValue   int32
	done        bool
}

//newLocal 构建本地处理服务
func newLocal(path string) *local {
	l := &local{
		FPS:       cmap.New(8),
		path:      path,
		lockValue: 0,
		init:      true,
		fpPath:    filepath.Join(path, ".fp"),
	}
	l.nfsChecker.Add(1)
	l.fsWatcher, _ = fsnotify.NewWatcher()
	l.fsWatcher.Add(path)
	go l.loopWatch()
	return l
}
func (l *local) loopWatch() {
	if l.fsWatcher == nil {
		return
	}
	for {
		select {
		case ev, ok := <-l.fsWatcher.Events:
			if !ok {
				return
			}
			fmt.Println("文件变化:", filepath.Base(ev.Name))
			if strings.HasPrefix(filepath.Base(ev.Name), ".") {
				continue
			}
			if ev.Op&fsnotify.Create == fsnotify.Create ||
				ev.Op&fsnotify.Write == fsnotify.Write ||
				ev.Op&fsnotify.Remove == fsnotify.Remove ||
				ev.Op&fsnotify.Rename == fsnotify.Rename ||
				ev.Op&fsnotify.Chmod == fsnotify.Chmod {

				fmt.Println("文件变化:", ev.Name)
				if err := l.check(); err != nil {
					hydra.G.Log().Debug("check.err", err)
				}
			}

		}
	}
}

//Update 更新配置数据
func (l *local) Update(currentAddr string) {
	needCheck := l.currentAddr != currentAddr
	l.currentAddr = currentAddr
	l.fpPath = filepath.Join(l.path, ".fp")
	if !needCheck {
		return
	}
	if err := l.check(); err != nil {
		hydra.G.Log().Debug("check.err", err)
	}
	return
}

//Merge 合并到本地列表
func (l *local) Merge(list eFileFPLists) (reports eFileFPLists, download eFileFPLists, err error) {
	reports = make(eFileFPLists, 10)
	download = make(eFileFPLists, 10)
	for _, fp := range list {
		nlk, ok := l.FPS.Get(fp.Path)
		if !ok {
			download[fp.Path] = fp
			continue
		}
		lk := nlk.(*eFileFP)
		v0 := fp.MergeHosts(l.currentAddr)
		v1 := lk.MergeHosts(fp.Hosts...)
		v2 := fp.MergeHosts(lk.Hosts...)
		if v0 || v1 || v2 {
			reports[fp.Path] = fp
		}
	}
	if len(reports) > 0 {
		err = l.FPWrite(l.FPS.Items())
	}
	return reports, download, err
}

//Open 读取文件
func (l *local) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.path, name))
}

//Close 将缓存数据写入本地文件
func (l *local) Close() error {
	l.done = true
	l.once.Do(func() {
		if l.fsWatcher != nil {
			l.fsWatcher.Close()
		}
		l.FPS.Clear()
	})
	return nil
}

//check 处理本地文件与指纹不一致，以文件为准
func (l *local) check() error {
	//只允许一个协程执行检查任务
	if !atomic.CompareAndSwapInt32(&l.lockValue, 0, 1) {
		return nil
	}
	//添加等待任务，完成后结束任务，并释放任务数
	l.nfsChecker.Add(1)
	defer atomic.CompareAndSwapInt32(&l.lockValue, 1, 0)
	defer l.nfsChecker.Done()
	defer func() {
		if l.init {
			l.init = false
			l.nfsChecker.Done()
		}
	}()

	//读取本地指纹
	fps, err := l.FPRead()
	if err != nil {
		return err
	}

	//获取本地文件列表
	lst, err := l.FList(l.path)
	if err != nil {
		return err
	}
	//处理不一致数据
	for _, entity := range lst {
		fp := &eFileFP{
			Path:    entity.Path,
			Size:    entity.Size,
			ModTime: entity.ModTime,
			Hosts:   []string{l.currentAddr},
		}
		if v, ok := fps[entity.Path]; ok {
			fp.MergeHosts(v.Hosts...)
		}
		l.FPS.Set(entity.Path, fp)
	}

	lstMap := lst.GetMap()
	for _, fp := range fps {
		if _, ok := lstMap[fp.Path]; !ok {
			delete(fps, fp.Path)
			l.FPS.Remove(fp.Path)
		}
	}

	//更新数据
	return l.FPWrite(l.FPS.Items())
}
