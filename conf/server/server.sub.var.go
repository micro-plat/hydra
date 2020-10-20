package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type varSub struct {
	cnf  conf.IVarConf
	rlog *VarLoader
}

func newVarSub(cnf conf.IVarConf) *varSub {
	s := &varSub{cnf: cnf}
	s.rlog = GetVarLoader(cnf, func(cnf conf.IVarConf) (interface{}, error) {
		return rlog.GetConf(cnf)
	})
	return s
}

//GetRLogConf 获取远程日志配置
func (s *varSub) GetRLogConf() (*rlog.Layout, error) {
	rlogObj, err := s.rlog.GetConf()
	if err != nil {
		return nil, err
	}
	return rlogObj.(*rlog.Layout), nil
}
