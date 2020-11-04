package vars

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
)

//VarLoader 配置加载器
type VarLoader struct {
	obj  interface{}
	once sync.Once
	cnf  conf.IVarConf
	err  error
	f    func(cnf conf.IVarConf) (interface{}, error)
}

//GetVarLoader 获取配置加载器
func GetVarLoader(cnf conf.IVarConf, f func(cnf conf.IVarConf) (interface{}, error)) *VarLoader {
	return &VarLoader{
		cnf: cnf,
		f:   f,
	}
}

//GetConf 获取配置信息
func (l *VarLoader) GetConf() (interface{}, error) {
	l.once.Do(func() {
		l.obj, l.err = l.f(l.cnf)
	})
	return l.obj, l.err
}
