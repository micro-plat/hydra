package server

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
)

//Loader 配置加载器
type Loader struct {
	obj  interface{}
	once sync.Once
	cnf  conf.IMainConf
	f    func(cnf conf.IMainConf) interface{}
}

//GetLoader 获取配置加载器
func GetLoader(cnf conf.IMainConf, f func(cnf conf.IMainConf) interface{}) *Loader {
	return &Loader{
		cnf: cnf,
		f:   f,
	}
}

//GetConf 获取配置信息
func (l *Loader) GetConf() interface{} {
	l.once.Do(func() {
		l.obj = l.f(l.cnf)
	})
	return l.obj
}
