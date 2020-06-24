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
	s.rlog = GetVarLoader(cnf, rlog.ConfHandler(rlog.GetConf).Handle)
	return s
}

//GetRLogConf 获取远程日志配置
func (s *varSub) GetRLogConf() *rlog.Layout {
	return s.rlog.GetConf().(*rlog.Layout)
}
