package wvalue

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//MultiValueWatcher 配置监控服务
type MultiValueWatcher struct {
	Watchers   []*SingleValueWatcher //@fix 方便测试
	notifyChan chan *watcher.ValueChangeArgs
}

//NewMultiValueWatcher 监控服务器变化
func NewMultiValueWatcher(rgst registry.IRegistry, path []string, logger logger.ILogging) (w *MultiValueWatcher, err error) {
	w = &MultiValueWatcher{
		notifyChan: make(chan *watcher.ValueChangeArgs, 10),
	}
	w.Watchers = make([]*SingleValueWatcher, 0, len(path))
	for _, p := range path {
		watcher := NewSingleValueWatcher(rgst, p, logger)
		w.Watchers = append(w.Watchers, watcher)
	}
	return
}

//Start 开始监听所有节点变化
func (c *MultiValueWatcher) Start() (chan *watcher.ValueChangeArgs, error) {
	for _, watcher := range c.Watchers {
		watcher.notifyChan = c.notifyChan
		if _, err := watcher.Start(); err != nil {
			return nil, err
		}
	}
	return c.notifyChan, nil
}

//Close 关闭监控器
func (c *MultiValueWatcher) Close() {
	for _, wacher := range c.Watchers {
		wacher.Close()
	}
}

type valueFactory struct{}

func (f *valueFactory) Create(rgst registry.IRegistry, path []string, logger logger.ILogging) (watcher.IValueWatcher, error) {
	return NewMultiValueWatcher(rgst, path, logger)
}

func init() {
	watcher.RegisterWatcher(&valueFactory{})
}
