package global

import "github.com/micro-plat/lib4go/tgo"

var modules []*tgo.Module = make([]*tgo.Module, 0, 10)

//AddTGOModules 添加TGO模块
func AddTGOModules(modules ...*tgo.Module) {
	modules = append(modules, modules...)
}

//GetTGOModules 获取TGO模块
func GetTGOModules() []*tgo.Module {
	return modules
}
