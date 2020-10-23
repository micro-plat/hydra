package wchild

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//MultiChildWatcher 配置监控服务
type MultiChildWatcher struct {
	Watchers   []*ChildWatcher
	notifyChan chan *watcher.ChildChangeArgs
}

//NewMultiChildWatcher 监控服务器变化
func NewMultiChildWatcher(rgst registry.IRegistry, path []string, logger logger.ILogging) (w *MultiChildWatcher, err error) {
	w = &MultiChildWatcher{
		notifyChan: make(chan *watcher.ChildChangeArgs, 100),
	}
	w.Watchers = make([]*ChildWatcher, 0, len(path))
	for _, p := range path {
		watcher := NewChildWatcher(rgst, p, logger)
		w.Watchers = append(w.Watchers, watcher)
	}
	return
}

//Start 开始监听所有节点变化
func (c *MultiChildWatcher) Start() (chan *watcher.ChildChangeArgs, error) {
	for _, watcher := range c.Watchers {
		watcher.notifyChan = c.notifyChan
		if _, err := watcher.Start(); err != nil {
			return nil, err
		}
	}
	return c.notifyChan, nil
}

//Close 关闭监控器
func (c *MultiChildWatcher) Close() {
	for _, wacher := range c.Watchers {
		wacher.Close()
	}
}

type childFactory struct{}

func (f *childFactory) Create(rgst registry.IRegistry, path []string, logger logger.ILogging) (watcher.IChildWatcher, error) {
	return NewMultiChildWatcher(rgst, path, logger)
}

func init() {
	watcher.RegisterWatcher(&childFactory{})
}
