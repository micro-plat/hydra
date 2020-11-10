package watcher

import (
	"strings"

	"github.com/micro-plat/hydra/registry"
)

//IChildWatcher 注册中心节点监控
type IChildWatcher interface {
	Start() (chan *ChildChangeArgs, error)
	Close()
}

//IValueWatcher 注册中心节点监控
type IValueWatcher interface {
	Start() (chan *ValueChangeArgs, error)
	Close()
}

const (
	//ADD 新增节点
	ADD = iota + 1
	//CHANGE 节点变更
	CHANGE
	//DEL 删除节点
	DEL
)

//ChildChangeArgs 子节点变化通知事件
type ChildChangeArgs struct {
	Registry registry.IRegistry
	Deep     int
	Name     string
	Parent   string
	Children []string
	Version  int32
	OP       int
}

//NewCArgsByChange 构建子节点变化参数
func NewCArgsByChange(op int, deep int, parent string, children []string, v int32, r registry.IRegistry) *ChildChangeArgs {
	names := strings.Split(strings.Trim(parent, "/"), "/")
	return &ChildChangeArgs{OP: op,
		Registry: r,
		Parent:   parent,
		Version:  v,
		Children: children,
		Deep:     deep,
		Name:     names[len(names)-1],
	}
}

//ValueChangeArgs 节点变化信息
type ValueChangeArgs struct {
	Registry registry.IRegistry
	Path     string
	Content  []byte
	Version  int32
	OP       int
}
