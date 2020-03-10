package registry

import (
	"strings"
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
