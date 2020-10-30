package server

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
)

//Loader 配置加载器
type Loader struct {
	obj  interface{}
	once sync.Once
	cnf  conf.IServerConf
	err  error
	f    func(cnf conf.IServerConf) (interface{}, error)
}

//GetLoader 获取配置加载器
func GetLoader(cnf conf.IServerConf, f func(cnf conf.IServerConf) (interface{}, error)) *Loader {
	return &Loader{
		cnf: cnf,
		f:   f,
	}
}

//GetConf 获取配置信息
func (l *Loader) GetConf() (interface{}, error) {
	l.once.Do(func() {
		l.obj, l.err = l.f(l.cnf)
	})

	return l.obj, l.err
}
