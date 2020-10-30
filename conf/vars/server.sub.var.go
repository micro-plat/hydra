package vars

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type VarSub struct {
	cnf  conf.IVarConf
	rlog *VarLoader
}

func NewVarSub(cnf conf.IVarConf) *VarSub {
	s := &VarSub{cnf: cnf}
	s.rlog = GetVarLoader(cnf, func(cnf conf.IVarConf) (interface{}, error) {
		return rlog.GetConf(cnf)
	})
	return s
}

//GetRLogConf 获取远程日志配置
func (s *VarSub) GetRLogConf() (*rlog.Layout, error) {
	rlogObj, err := s.rlog.GetConf()
	if err != nil {
		return nil, err
	}
	return rlogObj.(*rlog.Layout), nil
}
