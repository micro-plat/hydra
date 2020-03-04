package watcher

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//MultiChildWatcher 配置监控服务
type MultiChildWatcher struct {
	watchers   []*ChildWatcher
	notifyChan chan *registry.ChildChangeArgs
}

//NewMultiChildWatcher 监控服务器变化
func NewMultiChildWatcher(rgst registry.IRegistry, path []string, logger logger.ILogging) (w *MultiChildWatcher, err error) {
	w = &MultiChildWatcher{
		notifyChan: make(chan *registry.ChildChangeArgs, 10),
	}
	w.watchers = make([]*ChildWatcher, 0, len(path))
	for _, path := range path {
		watcher := NewChildWatcher(rgst, path, logger)
		w.watchers = append(w.watchers, watcher)
	}
	return
}

//Start 开始监听所有节点变化
func (c *MultiChildWatcher) Start() (chan *registry.ChildChangeArgs, error) {
	for _, watcher := range c.watchers {
		watcher.notifyChan = c.notifyChan
		if _, err := watcher.Start(); err != nil {
			return nil, err
		}
	}
	return c.notifyChan, nil
}

//Close 关闭监控器
func (c *MultiChildWatcher) Close() {
	for _, wacher := range c.watchers {
		wacher.Close()
	}
}

type childFactory struct{}

func (f *childFactory) Create(rgst registry.IRegistry, path []string, logger logger.ILogging) (registry.IChildWatcher, error) {
	return NewMultiChildWatcher(rgst, path, logger)
}

func init() {
	registry.RegisterWatcher(&childFactory{})
}
