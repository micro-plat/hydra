package watcher

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//ContentChangeArgs 值变化通知事件
type ContentChangeArgs struct {
	Path    string
	Content []byte
	Version int32
	OP      int
}

//Watcher 监控器
type Watcher struct {
	path       string
	timeSpan   time.Duration
	changed    int32
	notifyChan chan *ContentChangeArgs
	logger     *logger.Logger
	registry   registry.IRegistry
	mu         sync.Mutex
	done       bool
	closeChan  chan struct{}
}

//NewWatcher 初始化监控
func NewWatcher(path string, timeSpan time.Duration, registry registry.IRegistry, logger *logger.Logger) *Watcher {
	return &Watcher{
		path:       path,
		timeSpan:   timeSpan,
		registry:   registry,
		logger:     logger,
		notifyChan: make(chan *ContentChangeArgs, 1),
		closeChan:  make(chan struct{}),
	}
}

//Start 监控配置项变化，当发生错误时持续监控节点变化，只有明确节点不存在时才会通知关闭
func (w *Watcher) Start() (c chan *ContentChangeArgs, err error) {
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

func (w *Watcher) watch() (err error) {
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
func (w *Watcher) Close() {
	w.done = true
	close(w.closeChan)
}
func (w *Watcher) notifyDeleted() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v := atomic.LoadInt32(&w.changed); v > 0 && atomic.CompareAndSwapInt32(&w.changed, v, 0) {
		updater := &ContentChangeArgs{OP: DEL, Path: w.path}
		w.notifyChan <- updater

	}
}
func (w *Watcher) notifyChanged(content []byte, version int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := ADD
	if atomic.LoadInt32(&w.changed) > 0 {
		op = CHANGE
	}
	atomic.AddInt32(&w.changed, 1)
	updater := &ContentChangeArgs{OP: op, Path: w.path, Version: version, Content: content}
	w.notifyChan <- updater
	return
}
