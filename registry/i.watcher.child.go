package registry

import "strings"

//ChildrenChangeArgs 子节点变化通知事件
type ChildrenChangeArgs struct {
	Deep     int
	Name     string
	Parent   string
	Children []string
	Version  int32
	OP       int
}

//NewCArgsByChange 构建子节点变化参数
func NewCArgsByChange(op int, deep int, parent string, chilren []string, v int32) *ChildrenChangeArgs {
	names := strings.Split(strings.Trim(parent, "/"), "/")
	return &ChildrenChangeArgs{OP: op,
		Parent:   parent,
		Version:  v,
		Children: chilren,
		Deep:     deep,
		Name:     names[len(names)-1],
	}
}
