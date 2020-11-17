package global

import (
	"github.com/micro-plat/lib4go/tgo"
)

var modules []*tgo.Module = make([]*tgo.Module, 0, 10)

//AddTGOModules 添加TGO模块
func AddTGOModules(m ...*tgo.Module) {
	modules = append(modules, m...)
}

//GetTGOModules 获取TGO模块
func GetTGOModules() []*tgo.Module {
	return modules
}
