package watcher

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//SingleValueWatcher 监控器
type SingleValueWatcher struct {
	path       string
	timeSpan   time.Duration
	changed    int32
	notifyChan chan *registry.NodeChangeArgs
	logger     logger.ILogging
	registry   registry.IRegistry
	mu         sync.Mutex
	done       bool
	closeChan  chan struct{}
}

//NewSingleValueWatcher 监控节点值发生变化
func NewSingleValueWatcher(r registry.IRegistry, path string, logger logger.ILogging) *SingleValueWatcher {
	return &SingleValueWatcher{
		path:       path,
		timeSpan:   time.Second,
		registry:   r,
		logger:     logger,
		notifyChan: make(chan *registry.NodeChangeArgs, 1),
		closeChan:  make(chan struct{}),
	}
}

//Start 监控配置项变化，当发生错误时持续监控节点变化，只有明确节点不存在时才会通知关闭
func (w *SingleValueWatcher) Start() (c chan *registry.NodeChangeArgs, err error) {
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

func (w *SingleValueWatcher) watch() (err error) {
LOOP:
	exists, _ := w.registry.Exists(w.path)
	for !exists {
		select {
		case <-time.After(w.timeSpan):
			if w.done {
				return nil
			}
			exists, err = w.registry.Exists(w.path)
			if !exists && err == nil {
				w.notifyDeleted()
			}
		}
	}

	//获取节点值
	data, version, err := w.registry.GetValue(w.path)
	if err != nil {
		w.logger.Debugf("获取节点值失败：%s(err:%v)", w.path, err)
		time.Sleep(time.Second)
		goto LOOP
	}
	w.notifyChanged(data, version)
	dataChan, err := w.registry.WatchValue(w.path)
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
			w.notifyChanged(content.GetValue())
			//继续监控值变化
			dataChan, err = w.registry.WatchValue(w.path)
			if err != nil {
				goto LOOP
			}
		}
	}
}

//Close 关闭监控
func (w *SingleValueWatcher) Close() {
	w.done = true
	close(w.closeChan)
}
func (w *SingleValueWatcher) notifyDeleted() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v := atomic.LoadInt32(&w.changed); v > 0 && atomic.CompareAndSwapInt32(&w.changed, v, 0) {
		updater := &registry.NodeChangeArgs{OP: registry.DEL, Path: w.path, Registry: w.registry}
		w.notifyChan <- updater

	}
}
func (w *SingleValueWatcher) notifyChanged(content []byte, version int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := registry.ADD
	if atomic.LoadInt32(&w.changed) > 0 {
		op = registry.CHANGE
	}
	atomic.AddInt32(&w.changed, 1)
	updater := &registry.NodeChangeArgs{OP: op, Path: w.path, Version: version, Content: content, Registry: w.registry}
	w.notifyChan <- updater
	return
}
