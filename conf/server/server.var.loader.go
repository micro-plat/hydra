package server

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
)

//VarLoader 配置加载器
type VarLoader struct {
	obj  interface{}
	once sync.Once
	cnf  conf.IVarConf
	f    func(cnf conf.IVarConf) interface{}
}

//GetVarLoader 获取配置加载器
func GetVarLoader(cnf conf.IVarConf, f func(cnf conf.IVarConf) interface{}) *VarLoader {
	return &VarLoader{
		cnf: cnf,
		f:   f,
	}
}

//GetConf 获取配置信息
func (l *VarLoader) GetConf() interface{} {
	l.once.Do(func() {
		l.obj = l.f(l.cnf)
	})
	return l.obj
}
