package conf

import (
	"errors"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//IServerConf 服务器配置信息
type IServerConf interface {
	GetMainConf() IMainConf
	GetVarConf() IVarConf
}

//ServerConf 服务器配置信息
type ServerConf struct {
	mainConf IMainConf
	varConf  IVarConf
}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, rgst registry.IRegistry) (s *ServerConf, err error) {

	platName, sysName, serverType, clusterName := Split(mainConfpath)
	s = &ServerConf{}
	s.mainConf, err = NewMainConf(platName, sysName, serverType, clusterName, rgst)
	if err != nil {
		return nil, err
	}
	s.varConf, err = NewVarConf(s.mainConf.GetVarPath(), rgst)
	if err != nil {
		return nil, err
	}
	return s, nil

}

//GetMainConf 获取服务器主配置
func (s *ServerConf) GetMainConf() IMainConf {
	return s.mainConf
}

//GetVarConf 获取变量配置
func (s *ServerConf) GetVarConf() IVarConf {
	return s.varConf
}
