package watcher

import (
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

const (
	//ADD 新增节点
	ADD = iota + 1
	//CHANGE 节点变更
	CHANGE
	//DEL 删除节点
	DEL
)

//ConfWatcher 配置监控服务
type ConfWatcher struct {
	paths      []string
	watchers   []*Watcher
	timeSpan   time.Duration
	logger     *logger.Logger
	notifyChan chan *ContentChangeArgs
	registry   registry.IRegistry
}

//NewConfWatcher 初始化服务器监控程序
func NewConfWatcher(platName string, systemName string, serverTypes []string, clusterName string, rgst registry.IRegistry, logger *logger.Logger) (w *ConfWatcher, err error) {
	w = &ConfWatcher{
		timeSpan:   time.Second,
		registry:   rgst,
		logger:     logger,
		notifyChan: make(chan *ContentChangeArgs, 10),
	}
	w.paths = make([]string, 0, len(serverTypes))
	for _, tp := range serverTypes {
		w.paths = append(w.paths, registry.Join("/", platName, systemName, tp, clusterName, "conf"))
	}

	w.watchers = make([]*Watcher, 0, len(w.paths))
	for _, path := range w.paths {
		watcher := NewWatcher(path, w.timeSpan, rgst, logger)
		w.watchers = append(w.watchers, watcher)
	}
	return
}

//Notify 服务器变更通知
func (c *ConfWatcher) Notify() (chan *ContentChangeArgs, error) {
	for _, watcher := range c.watchers {
		watcher.notifyChan = c.notifyChan
		if _, err := watcher.Start(); err != nil {
			return nil, err
		}
	}
	return c.notifyChan, nil
}

//Close 关闭监控器
func (c *ConfWatcher) Close() {
	for _, wacher := range c.watchers {
		wacher.Close()
	}
}
