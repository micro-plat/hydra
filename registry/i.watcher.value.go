package registry

import (
	"fmt"
	"strings"

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

const (
	//WatchValue 监控节点值变化
	WatchValue = "VALUE"

	//WatchChild 监控子节点变化
	WatchChild = "CHILD"
)

//IValueWatcher 注册中心节点监控
type IValueWatcher interface {
	Start() (chan *ValueChangeArgs, error)
	Close()
}

//ValueChangeArgs 节点变化信息
type ValueChangeArgs struct {
	Registry IRegistry
	Path     string
	Content  []byte
	Version  int32
	OP       int
}

//IsConf 是否是conf根节点或conf的子节点
func (n *ValueChangeArgs) IsConf() bool {
	return strings.HasSuffix(n.Path, Join(n.Registry.GetSeparator(), "conf")) ||
		strings.Contains(n.Path, Join(n.Registry.GetSeparator(), "conf", n.Registry.GetSeparator()))
}

//IsVarRoot 是否是var跟节点或var的子节点
func (n *ValueChangeArgs) IsVarRoot() bool {
	return strings.HasSuffix(n.Path, Join(n.Registry.GetSeparator(), "var")) ||
		strings.Contains(n.Path, Join(n.Registry.GetSeparator(), "var", n.Registry.GetSeparator()))
}

//IValueWatcherFactory watcher生成器
type IValueWatcherFactory interface {
	Create(r IRegistry, path []string, l logger.ILogging) (IValueWatcher, error)
}

var valueWatchers = make(map[string]IValueWatcherFactory)

//RegisterValueWatcher 注册配置文件适配器
func RegisterValueWatcher(name string, f IValueWatcherFactory) {
	if f == nil {
		panic("registry: 注册watcher factory不能为空")
	}
	if _, ok := valueWatchers[name]; ok {
		panic("registry: 重复注册watcher factory " + name)
	}
	valueWatchers[name] = f
}

//NewValueWatcher 构建值监控,监控指定路径的值变化
func NewValueWatcher(registryAddr string, path []string, l logger.ILogging) (IValueWatcher, error) {
	r, err := NewRegistry(registryAddr, l)
	if err != nil {
		return nil, err
	}
	return valueWatchers[WatchValue].Create(r, path, l)
}

//NewValueWatcherByRegistry 根据注册中心构建监控对象，监控指定路径的值变化
func NewValueWatcherByRegistry(registry IRegistry, path []string, l logger.ILogging) (IValueWatcher, error) {
	return valueWatchers[WatchValue].Create(registry, path, l)
}

//NewValueWatcherByServers 根据服务器信息，监控所有服务器节点发生变化
func NewValueWatcherByServers(registry IRegistry, platName string, systemName string, serverTypes []string, clusterName string, l logger.ILogging) (IValueWatcher, error) {
	if platName == "" || systemName == "" || len(serverTypes) == 0 || clusterName == "" {
		return nil, fmt.Errorf("指定的平台名称:%s，系统名称:%s，服务类型:%v，集群名称:%s,不能为空", platName, systemName, serverTypes, clusterName)
	}
	path := make([]string, 0, len(serverTypes))
	for _, s := range serverTypes {
		path = append(path, Join(platName, systemName, s, clusterName, "conf"))
	}
	return NewValueWatcherByRegistry(registry, path, l)
}
