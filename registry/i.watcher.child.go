package registry

import (
	"strings"

	"github.com/micro-plat/lib4go/logger"
)

//ChildChangeArgs 子节点变化通知事件
type ChildChangeArgs struct {
	Registry IRegistry
	Deep     int
	Name     string
	Parent   string
	Children []string
	Version  int32
	OP       int
}

//NewCArgsByChange 构建子节点变化参数
func NewCArgsByChange(op int, deep int, parent string, chilren []string, v int32, r IRegistry) *ChildChangeArgs {
	names := strings.Split(strings.Trim(parent, r.GetSeparator()), r.GetSeparator())
	return &ChildChangeArgs{OP: op,
		Registry: r,
		Parent:   parent,
		Version:  v,
		Children: chilren,
		Deep:     deep,
		Name:     names[len(names)-1],
	}
}

//IChildWatcher 注册中心节点监控
type IChildWatcher interface {
	Start() (chan *ChildChangeArgs, error)
	Close()
}

//IChildWatcherFactory watcher生成器
type IChildWatcherFactory interface {
	Create(r IRegistry, path []string, l logger.ILogging) (IChildWatcher, error)
}

var childWatchers = make(map[string]IChildWatcherFactory)

//RegisterChildWatcher 注册配置文件适配器
func RegisterChildWatcher(name string, f IChildWatcherFactory) {
	if f == nil {
		panic("registry: 注册child watcher factory不能为空")
	}
	if _, ok := childWatchers[name]; ok {
		panic("registry: 重复注册child watcher factory " + name)
	}
	childWatchers[name] = f
}

//NewChildWatcher 构建值监控,监控指定路径的值变化
func NewChildWatcher(registryAddr string, path []string, l logger.ILogging) (IChildWatcher, error) {
	r, err := NewRegistry(registryAddr, l)
	if err != nil {
		return nil, err
	}
	return childWatchers[WatchValue].Create(r, path, l)
}
