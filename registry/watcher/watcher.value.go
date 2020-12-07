package watcher

import (
	"fmt"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//NewValueWatcher 构建值监控,监控指定路径的值变化
func NewValueWatcher(registryAddr string, path []string, l logger.ILogging) (IValueWatcher, error) {
	r, err := registry.GetRegistry(registryAddr, l)
	if err != nil {
		return nil, err
	}
	return NewValueWatcherByRegistry(r, path, l)
}

//NewValueWatcherByRegistry 根据注册中心构建监控对象，监控指定路径的值变化
func NewValueWatcherByRegistry(r registry.IRegistry, path []string, l logger.ILogging) (IValueWatcher, error) {
	if valueWatchFactory == nil {
		panic(fmt.Errorf("未提供节点值监控的工厂实现对象IValueWatcherFactory"))
	}
	return valueWatchFactory.Create(r, path, l)
}

//NewValueWatcherByServers 根据服务器信息，监控所有服务器节点发生变化
func NewValueWatcherByServers(r registry.IRegistry, platName string, systemName string, serverTypes []string, clusterName string, l logger.ILogging) (IValueWatcher, error) {
	if platName == "" || systemName == "" || len(serverTypes) == 0 || clusterName == "" {
		panic(fmt.Errorf("指定的平台名称:%s，系统名称:%s，服务类型:%v，集群名称:%s,不能为空", platName, systemName, serverTypes, clusterName))
	}
	path := make([]string, 0, len(serverTypes))
	for _, s := range serverTypes {
		path = append(path, registry.Join(platName, systemName, s, clusterName, "conf"))
	}
	return NewValueWatcherByRegistry(r, path, l)
}
