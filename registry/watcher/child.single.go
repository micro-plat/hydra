package watcher

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//ChildWatcher 监控器
type ChildWatcher struct {
	path       string
	timeSpan   time.Duration
	deep       int
	changedCnt int32
	notifyChan chan *registry.ChildChangeArgs
	logger     logger.ILogging
	registry   registry.IRegistry
	watchers   map[string]*ChildWatcher
	mu         sync.Mutex
	done       bool
	closeChan  chan struct{}
}

//NewChildWatcher 初始化子节点监控
func NewChildWatcher(registry registry.IRegistry, path string, logger logger.ILogging) *ChildWatcher {
	return NewChildWatcherByDeep(path, 2, registry, logger)
}

//NewChildWatcherByDeep 初始化监控
func NewChildWatcherByDeep(path string, deep int, r registry.IRegistry, logger logger.ILogging) *ChildWatcher {
	return &ChildWatcher{
		path:       path,
		deep:       deep,
		timeSpan:   time.Second,
		registry:   r,
		logger:     logger,
		watchers:   make(map[string]*ChildWatcher),
		notifyChan: make(chan *registry.ChildChangeArgs, 1),
		closeChan:  make(chan struct{}),
	}
}

//Start 监控配置项变化，当发生错误时持续监控节点变化，只有明确节点不存在时才会通知关闭
func (w *ChildWatcher) Start() (c chan *registry.ChildChangeArgs, err error) {
	errChan := make(chan error, 1)
	go func() {
		err := w.watch(w.path)
		if err != nil {
			errChan <- err
		}
	}()
	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Microsecond * 500):
		return w.notifyChan, nil
	}
}
func (w *ChildWatcher) watch(path string) (err error) {
LOOP:
	exists, _ := w.registry.Exists(path)
	for !exists {
		select {
		case <-time.After(w.timeSpan):
			if w.done {
				return nil
			}
			exists, err = w.registry.Exists(path)
			if !exists && err == nil {
				w.deleted()
			}
		}
	}
	//获取节点值
	data, version, err := w.registry.GetChildren(path)
	if err != nil {
		w.logger.Debugf("获取节点值失败：%s(err:%v)", path, err)
		time.Sleep(time.Second)
		goto LOOP
	}
	w.changed(data, version)
	dataChan, err := w.registry.WatchChildren(path)
	if err != nil {
		goto LOOP
	}

	for {
		select {
		case <-w.closeChan:
			return nil
		case content, ok := <-dataChan:
			if w.done || !ok {
				return nil
			}
			if err = content.GetError(); err != nil {
				goto LOOP
			}

			if b, _ := w.registry.Exists(path); !b {
				w.deleted()
				goto LOOP
			}

			data, version := content.GetValue()
			w.changed(data, version)
			//继续监控值变化
			dataChan, err = w.registry.WatchChildren(path)
			if err != nil {
				goto LOOP
			}
		}
	}
}

//Close 关闭监控
func (w *ChildWatcher) Close() {
	for _, watcher := range w.watchers {
		watcher.Close()
	}
	w.done = true
	close(w.closeChan)
}

//deleted 节点删除
func (w *ChildWatcher) deleted() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v := atomic.LoadInt32(&w.changedCnt); v > 0 && atomic.CompareAndSwapInt32(&w.changedCnt, v, 0) {
		updater := registry.NewCArgsByChange(registry.DEL, w.deep, w.path, nil, 0, w.registry)
		w.notify(updater)
	}
}

//changed 子节点发生变化
func (w *ChildWatcher) changed(children []string, version int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := registry.ADD
	if atomic.LoadInt32(&w.changedCnt) > 0 {
		op = registry.CHANGE
	}
	atomic.AddInt32(&w.changedCnt, 1)
	updater := registry.NewCArgsByChange(op, w.deep, w.path, children, version, w.registry)
	w.notify(updater)
	return
}
func (w *ChildWatcher) notify(a *registry.ChildChangeArgs) {
	if a.Deep == 1 && a.OP != registry.DEL {
		w.notifyChan <- a
		return
	}
	for _, path := range a.Children {
		switch a.OP {
		case registry.ADD, registry.CHANGE:
			w.changeChilrenWatcher(path)
		case registry.DEL:
			w.notifyChan <- a
			w.delChilrenWatcher(path)
		}
	}
}

func (w *ChildWatcher) delChilrenWatcher(path string) {
	if w, ok := w.watchers[path]; ok {
		w.Close()
		delete(w.watchers, path)
	}

}

func (w *ChildWatcher) changeChilrenWatcher(path string) {

	if _, ok := w.watchers[path]; ok {
		return
	}
	watcher := NewChildWatcherByDeep(registry.Join(w.path, path), w.deep-1, w.registry, w.logger)
	ch, err := watcher.Start()
	if err != nil {
		w.logger.Error(err)
		return
	}
	w.watchers[path] = watcher
	go func() {
		for {
			select {
			case <-w.closeChan:
				return
			case arg, ok := <-ch:
				if w.done || !ok {
					return
				}
				w.notify(arg)
			}
		}
	}()
}