package watcher

import (
	"fmt"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//NewChildWatcher 构建值监控,监控指定路径的值变化
func NewChildWatcher(registryAddr string, path []string, l logger.ILogging) (IChildWatcher, error) {
	if childWatcherFactory == nil {
		panic(fmt.Errorf("未提供子节点监控工厂实现对象IChildWatcherFactory"))
	}
	r, err := registry.NewRegistry(registryAddr, l)
	if err != nil {
		return nil, err
	}
	return childWatcherFactory.Create(r, path, l)
}

//NewChildWatcherByRegistry 构建值监控,监控指定路径的值变化
func NewChildWatcherByRegistry(r registry.IRegistry, path []string, l logger.ILogging) (IChildWatcher, error) {
	if childWatcherFactory == nil {
		panic(fmt.Errorf("未提供子节点监控工厂实现对象IChildWatcherFactory"))
	}
	return childWatcherFactory.Create(r, path, l)
}
