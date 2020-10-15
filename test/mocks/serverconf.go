package mocks

import (
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
)

//SConf 服务器配置
type SConf struct {
	creator.IConf
}

//NewConf 构建配置信息
func NewConf() *SConf {
	return &SConf{
		IConf: creator.New(),
	}
}

//Conf 配置
func (s *SConf) Conf() creator.IConf {
	return s.IConf
}

//GetAPIConf 获取API服务器配置
func (s *SConf) GetAPIConf() server.IServerConf {
	return s.GetConf("hydra", "apiserver", "api", "test")
}

//GetWebConf 获取web服务器配置
func (s *SConf) GetWebConf() server.IServerConf {
	return s.GetConf("hydra", "webserver", "web", "test")
}

//GetWSConf 获取API服务器配置
func (s *SConf) GetWSConf() server.IServerConf {
	return s.GetConf("hydra", "wsserver", "ws", "test")
}

//GetCronConf 获取cron服务器配置
func (s *SConf) GetCronConf() server.IServerConf {
	return s.GetConf("hydra", "cronserver", "cron", "test")
}

//GetMQCConf 获取mqc服务器配置
func (s *SConf) GetMQCConf() server.IServerConf {
	return s.GetConf("hydra", "mqcserver", "mqc", "test")
}

//GetConf 获取配置信息
func (s *SConf) GetConf(platName string, systemName string, serverType string, clusterName string) server.IServerConf {
	registryAddr := "lm://."
	if err := s.IConf.Pub(platName, systemName, clusterName, registryAddr, true); err != nil {
		panic(err)
	}
	r, err := registry.NewRegistry(registryAddr, global.Def.Log())
	if err != nil {
		panic(err)
	}
	path := registry.Join(platName, systemName, serverType, clusterName, "conf")
	conf, err := server.NewServerConf(path, r)
	if err != nil {
		panic(err)
	}
	return conf
}
