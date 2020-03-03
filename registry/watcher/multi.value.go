package watcher

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//MultiWatcher 配置监控服务
type MultiWatcher struct {
	watchers   []*SingleWatcher
	notifyChan chan *registry.NodeChangeArgs
}

//NewMultiWatcher 监控服务器变化
func NewMultiWatcher(rgst registry.IRegistry, path []string, logger logger.ILogging) (w *MultiWatcher, err error) {
	w = &MultiWatcher{
		notifyChan: make(chan *registry.NodeChangeArgs, 10),
	}
	w.watchers = make([]*SingleWatcher, 0, len(path))
	for _, path := range path {
		watcher := NewSingleWatcher(rgst, path, logger)
		w.watchers = append(w.watchers, watcher)
	}
	return
}

//Start 开始监听所有节点变化
func (c *MultiWatcher) Start() (chan *registry.NodeChangeArgs, error) {
	for _, watcher := range c.watchers {
		watcher.notifyChan = c.notifyChan
		if _, err := watcher.Start(); err != nil {
			return nil, err
		}
	}
	return c.notifyChan, nil
}

//Close 关闭监控器
func (c *MultiWatcher) Close() {
	for _, wacher := range c.watchers {
		wacher.Close()
	}
}

type factory struct{}

func (f *factory) Create(rgst registry.IRegistry, path []string, logger logger.ILogging) (registry.IWatcher, error) {
	return NewMultiWatcher(rgst, path, logger)
}

func init() {
	registry.RegisterWatcher(registry.WatchValue, &factory{})
}
