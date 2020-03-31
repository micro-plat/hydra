package config

import (
	"errors"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//IServerConf 集群节点配置
type IServerConf interface {
	GetMainConf() IMainConf
	GetVarConf() IVarConf
}

//ServerConf 服务器配置信息
type ServerConf struct {
}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, rgst registry.IRegistry) (s *ServerConf, err error) {
	s = &ServerConf{}
	s.ClusterConf, err = NewClusterConf(mainConfpath, rgst)
	if err != nil {
		return nil, err
	}
	s.VarConf, err = NewVarConf(mainConfpath, rgst)
	if err != nil {
		return nil, err
	}
	return s, nil

}
