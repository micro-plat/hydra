package registry

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"
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
	//WatchValue 监控值变化
	WatchValue = "VALUE"
)

//IWatcher 注册中心节点监控
type IWatcher interface {
	Start() (chan *NodeChangeArgs, error)
	Close()
}

//NodeChangeArgs 节点变化信息
type NodeChangeArgs struct {
	Registry IRegistry
	Path     string
	Content  []byte
	Version  int32
	OP       int
}

//IsConf 是否是conf根节点或conf的子节点
func (n *NodeChangeArgs) IsConf() bool {
	return strings.HasSuffix(n.Path, Join(n.Registry.GetSeparator(), "conf")) ||
		strings.Contains(n.Path, Join(n.Registry.GetSeparator(), "conf", n.Registry.GetSeparator()))
}

//IsVarRoot 是否是var跟节点或var的子节点
func (n *NodeChangeArgs) IsVarRoot() bool {
	return strings.HasSuffix(n.Path, Join(n.Registry.GetSeparator(), "var")) ||
		strings.Contains(n.Path, Join(n.Registry.GetSeparator(), "var", n.Registry.GetSeparator()))
}

//IWatcherFactory watcher生成器
type IWatcherFactory interface {
	Create(r IRegistry, path []string, l logger.ILogging) (IWatcher, error)
}

var watchersMap = cmap.New(2)
var watchers = make(map[string]IWatcherFactory)

//RegisterWatcher 注册配置文件适配器
func RegisterWatcher(name string, f IWatcherFactory) {
	if f == nil {
		panic("registry: 注册watcher factory不能为空")
	}
	if _, ok := watchers[name]; ok {
		panic("registry: 重复注册watcher factory " + name)
	}
	watchers[name] = f
}

//NewValueWatcher 构建值监控,监控指定路径的值变化
func NewValueWatcher(registryAddr string, path []string, l logger.ILogging) (IWatcher, error) {
	r, err := NewRegistry(registryAddr, l)
	if err != nil {
		return nil, err
	}
	return watchers[WatchValue].Create(r, path, l)
}

//NewValueWatcherByRegistry 根据注册中心构建监控对象，监控指定路径的值变化
func NewValueWatcherByRegistry(registry IRegistry, path []string, l logger.ILogging) (IWatcher, error) {
	return watchers[WatchValue].Create(registry, path, l)
}

//NewValueWatcherByServers 根据服务器信息，监控所有服务器节点发生变化
func NewValueWatcherByServers(registry IRegistry, platName string, systemName string, serverTypes []string, clusterName string, l logger.ILogging) (IWatcher, error) {
	if platName == "" || systemName == "" || len(serverTypes) == 0 || clusterName == "" {
		return nil, fmt.Errorf("指定的平台名称:%s，系统名称:%s，服务类型:%v，集群名称:%s,不能为空", platName, systemName, serverTypes, clusterName)
	}
	path := make([]string, 0, len(serverTypes))
	for _, s := range serverTypes {
		path = append(path, Join(platName, systemName, s, clusterName, "conf"))
	}
	path = append(path, Join(platName, "var"))
	return NewValueWatcherByRegistry(registry, path, l)
}
