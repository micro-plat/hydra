package watcher

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//IChildWatcherFactory watcher生成器
type IChildWatcherFactory interface {
	Create(r registry.IRegistry, path []string, l logger.ILogging) (IChildWatcher, error)
}

//IValueWatcherFactory watcher生成器
type IValueWatcherFactory interface {
	Create(r registry.IRegistry, path []string, l logger.ILogging) (IValueWatcher, error)
}

var childWatcherFactory IChildWatcherFactory
var valueWatchFactory IValueWatcherFactory

//RegisterWatcher 注册配置文件适配器
func RegisterWatcher(f interface{}) {
	switch v := f.(type) {
	case IValueWatcherFactory:
		valueWatchFactory = v
	case IChildWatcherFactory:
		childWatcherFactory = v
	default:
		panic("注册的watcher类型(非IValueWatcherFactory或IChildWatcherFactory)不支持")
	}
}
