package file

import (
	"os"
	"time"

	"sync"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//DirWatcher 文件夹监控
type DirWatcher struct {
	callback func()
	files    cmap.ConcurrentMap
	lastTime time.Time
	timeSpan time.Duration
	closeCh  chan struct{}
	done     bool
	mu       sync.Mutex
}

//NewDirWatcher 构建脚本监控文件
func NewDirWatcher(callback func(), timeSpan time.Duration) *DirWatcher {
	w := &DirWatcher{callback: callback, lastTime: time.Now(), timeSpan: timeSpan}
	w.closeCh = make(chan struct{}, 1)
	w.files = cmap.New(4)
	go w.watch()
	return w
}

//Append 添加监控文件
func (w *DirWatcher) Append(path string) (err error) {
	//dir := filepath.Dir(path)
	//w.files.SetIfAbsent(dir, dir)
	w.files.SetIfAbsent(path, struct{}{})
	return nil
}

func (w *DirWatcher) watch() {
	for {
		select {
		case <-w.closeCh:
			return
		case <-time.After(w.timeSpan):
			if w.done {
				return
			}
			if w.checkChange() {
				w.callback()
			}
		}
	}
}

//checkChange 检查文件夹最后修改时间
func (w *DirWatcher) checkChange() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	change := false
	w.files.IterCb(func(path string, v interface{}) bool {
		fileinfo, err := os.Stat(path)
		if err != nil {
			return true //继续检查下一个文件
		}
		if fileinfo.ModTime().Sub(w.lastTime) > 0 {
			w.lastTime = fileinfo.ModTime()
			change = true
			return false //当前文件发生变化，退出不再检查
		}
		return true //继续检查下一个文件
	})
	return change
}

//Close 关闭服务
func (w *DirWatcher) Close() {
	close(w.closeCh)
	w.done = true
}
