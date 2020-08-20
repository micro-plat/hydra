package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/apm"
)

type apmSub struct {
	cnf conf.IMainConf
	apm *Loader
}

func newapmSub(cnf conf.IMainConf) *apmSub {
	s := &apmSub{cnf: cnf}
	s.apm = GetLoader(cnf, apm.ConfHandler(apm.GetConf).Handle)
	return s
}

//GetAPIKeyConf 获取apikey配置
func (s *apmSub) GetAPMConf() *apm.APM {
	return s.apm.GetConf().(*apm.APM)
}
