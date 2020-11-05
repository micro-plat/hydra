package wvalue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//SingleValueWatcher 监控器
type SingleValueWatcher struct {
	path       string
	timeSpan   time.Duration
	changed    int32
	notifyChan chan *watcher.ValueChangeArgs
	logger     logger.ILogging
	registry   registry.IRegistry
	mu         sync.Mutex
	Done       bool          //@fix 方便测试
	CloseChan  chan struct{} //@fix 方便测试
}

//NewSingleValueWatcher 监控节点值发生变化
func NewSingleValueWatcher(r registry.IRegistry, path string, logger logger.ILogging) *SingleValueWatcher {
	return &SingleValueWatcher{
		path:       path,
		timeSpan:   time.Second,
		registry:   r,
		logger:     logger,
		notifyChan: make(chan *watcher.ValueChangeArgs, 1),
		CloseChan:  make(chan struct{}), //@fix 方便测试
	}
}

//Start 监控配置项变化，当发生错误时持续监控节点变化，只有明确节点不存在时才会通知关闭
func (w *SingleValueWatcher) Start() (c chan *watcher.ValueChangeArgs, err error) {
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
			if w.Done {
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
		case <-w.CloseChan:
			return nil
		case content, ok := <-dataChan: //程序阻塞,等待节点值变动通知
			if w.Done || !ok {
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
	w.Done = true
	close(w.CloseChan)
}
func (w *SingleValueWatcher) notifyDeleted() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v := atomic.LoadInt32(&w.changed); v > 0 && atomic.CompareAndSwapInt32(&w.changed, v, 0) {
		updater := &watcher.ValueChangeArgs{OP: watcher.DEL, Path: w.path, Registry: w.registry}
		w.notifyChan <- updater

	}
}
func (w *SingleValueWatcher) notifyChanged(content []byte, version int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := watcher.ADD
	if atomic.LoadInt32(&w.changed) > 0 {
		op = watcher.CHANGE
	}
	atomic.AddInt32(&w.changed, 1)
	updater := &watcher.ValueChangeArgs{OP: op, Path: w.path, Version: version, Content: content, Registry: w.registry}
	w.notifyChan <- updater
	return
}
