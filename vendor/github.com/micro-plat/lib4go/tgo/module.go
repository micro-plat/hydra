package tgo

import (
	"github.com/d5/tengo/v2"
	"github.com/micro-plat/lib4go/types"
)

//Module 供脚本使用的模块信息
type Module struct {
	name    string
	objects map[string]tengo.Object
}

//NewModule 构建模块
func NewModule(name ...string) *Module {
	return &Module{
		name:    types.GetStringByIndex(name, 0),
		objects: make(map[string]tengo.Object),
	}
}

//Add 添加函数
func (m *Module) Add(method string, f CallableFunc) *Module {
	m.objects[method] = &UserFunction{Name: method, Value: f}
	return m
}

//Objects 对象列表
func (m *Module) Objects() map[string]tengo.Object {
	return m.objects
}
