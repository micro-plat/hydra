package watcher

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//ChildrenChangeArgs 值变化通知事件
type ChildrenChangeArgs struct {
	Deep     int
	Name     string
	Parent   string
	Children []string
	Version  int32
	OP       int
}

func NewCArgsByDelete(deep int, parent string) *ChildrenChangeArgs {
	return NewCArgsByChange(DEL, deep, parent, nil, 0)
}
func NewCArgsByChange(op int, deep int, parent string, chilren []string, v int32) *ChildrenChangeArgs {
	names := strings.Split(strings.Trim(parent, "/"), "/")
	return &ChildrenChangeArgs{OP: op,
		Parent:   parent,
		Version:  v,
		Children: chilren,
		Deep:     deep,
		Name:     names[len(names)-1],
	}
}

//ChildrenWatcher 监控器
type ChildrenWatcher struct {
	path       string
	timeSpan   time.Duration
	deep       int
	changed    int32
	notifyChan chan *ChildrenChangeArgs
	logger     logger.ILogger
	registry   registry.IRegistry
	watchers   map[string]*ChildrenWatcher
	mu         sync.Mutex
	done       bool
	closeChan  chan struct{}
}

//NewChildrenWatcher 初始化监控
func NewChildrenWatcher(path string, deep int, timeSpan time.Duration, registry registry.IRegistry, logger logger.ILogger) *ChildrenWatcher {
	return &ChildrenWatcher{
		path:       path,
		deep:       deep,
		timeSpan:   timeSpan,
		registry:   registry,
		logger:     logger,
		watchers:   make(map[string]*ChildrenWatcher),
		notifyChan: make(chan *ChildrenChangeArgs, 1),
		closeChan:  make(chan struct{}),
	}
}

//Start 监控配置项变化，当发生错误时持续监控节点变化，只有明确节点不存在时才会通知关闭
func (w *ChildrenWatcher) Start() (c chan *ChildrenChangeArgs, err error) {
	errChan := make(chan error, 1)
	go func() {
		err := w.watch()
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
func (w *ChildrenWatcher) watch() error {
	w.watchOne(w.path)
	return nil

}
func (w *ChildrenWatcher) watchOne(path string) (err error) {
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
				w.Deleted()
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
	w.Changed(data, version)
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
				w.Deleted()
				goto LOOP
			}

			data, version := content.GetValue()
			w.Changed(data, version)
			//继续监控值变化
			dataChan, err = w.registry.WatchChildren(path)
			if err != nil {
				goto LOOP
			}
		}
	}
}

//Close 关闭监控
func (w *ChildrenWatcher) Close() {
	for _, watcher := range w.watchers {
		watcher.Close()
	}
	w.done = true
	close(w.closeChan)
}

//Deleted 节点删除
func (w *ChildrenWatcher) Deleted() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v := atomic.LoadInt32(&w.changed); v > 0 && atomic.CompareAndSwapInt32(&w.changed, v, 0) {
		updater := NewCArgsByDelete(w.deep, w.path)
		w.notify(updater)
	}
}

//Changed 子节点发生变化
func (w *ChildrenWatcher) Changed(children []string, version int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := ADD
	if atomic.LoadInt32(&w.changed) > 0 {
		op = CHANGE
	}
	atomic.AddInt32(&w.changed, 1)
	updater := NewCArgsByChange(op, w.deep, w.path, children, version)
	w.notify(updater)
	return
}
func (w *ChildrenWatcher) notify(a *ChildrenChangeArgs) {
	if a.Deep == 1 && a.OP != DEL {
		w.notifyChan <- a
		return
	}
	for _, path := range a.Children {
		switch a.OP {
		case ADD, CHANGE:
			w.changeChilrenWatcher(path)
		case DEL:
			w.notifyChan <- a
			w.delChilrenWatcher(path)
		}
	}
}

func (w *ChildrenWatcher) delChilrenWatcher(path string) {
	if w, ok := w.watchers[path]; ok {
		w.Close()
		delete(w.watchers, path)
	}

}

func (w *ChildrenWatcher) changeChilrenWatcher(path string) {

	if _, ok := w.watchers[path]; ok {
		return
	}
	watcher := NewChildrenWatcher(registry.Join(w.path, path), w.deep-1, w.timeSpan, w.registry, w.logger)
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
