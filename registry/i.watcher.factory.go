package registry

import "github.com/micro-plat/lib4go/logger"

//IChildWatcherFactory watcher生成器
type IChildWatcherFactory interface {
	Create(r IRegistry, path []string, l logger.ILogging) (IChildWatcher, error)
}

//IValueWatcherFactory watcher生成器
type IValueWatcherFactory interface {
	Create(r IRegistry, path []string, l logger.ILogging) (IValueWatcher, error)
}

var childWatcherFactory IChildWatcherFactory
var valueWatchFactory IValueWatcherFactory

//RegisterValueWatcher 注册配置文件适配器
func RegisterWatcher(f interface{}) {
	switch (v:=f.(type)){
	case IValueWatcherFactory:
		valueWatchFactory=v
	case IChildWatcherFactory:
		childWatcherFactory=v
	default:
		panic("注册的watcher类型不支持")
}
}
