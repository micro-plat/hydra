package mocks

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/services"
)

type service struct {
	API *services.ORouter
	Web *services.ORouter
	WS  *services.ORouter
	RPC *services.ORouter
}

//SConf 服务器配置
type SConf struct {
	creator.IConf
	PlatName    string
	ClusterName string
	Service     *service
}

//NewConf 构建配置信息
func NewConf() *SConf {
	c := &SConf{
		PlatName:    "hydra",
		ClusterName: "test",
		Service:     &service{},
	}
	//API  路由信息
	c.Service.API = services.NewORouter()

	//WEB web服务的路由信息
	c.Service.Web = services.NewORouter()

	//WS web socket路由信息
	c.Service.WS = services.NewORouter()

	//RPC rpc服务的路由信息
	c.Service.RPC = services.NewORouter()

	c.IConf = creator.New(c.getRouter)

	//处理iconf.load中，服务检查问题
	global.Def.ServerTypes = []string{http.API, http.Web, http.WS, cron.CRON}
	// hydra.WithServerTypes(http.API)
	return c
}

//Conf 配置
func (s *SConf) Conf() creator.IConf {
	return s.IConf
}

//GetAPIConf 获取API服务器配置
func (s *SConf) GetAPIConf() server.IServerConf {
	return s.GetConf(s.PlatName, "apiserver", "api", s.ClusterName)
}

//GetWebConf 获取web服务器配置
func (s *SConf) GetWebConf() server.IServerConf {
	return s.GetConf(s.PlatName, "webserver", "web", s.ClusterName)
}

//GetWSConf 获取API服务器配置
func (s *SConf) GetWSConf() server.IServerConf {
	return s.GetConf(s.PlatName, "wsserver", "ws", s.ClusterName)
}

//GetCronConf 获取cron服务器配置
func (s *SConf) GetCronConf() server.IServerConf {
	return s.GetConf(s.PlatName, "cronserver", "cron", s.ClusterName)
}

//GetMQCConf 获取mqc服务器配置
func (s *SConf) GetMQCConf() server.IServerConf {
	global.Def.ServerTypes = []string{http.API, http.Web, http.WS, cron.CRON, mqc.MQC}
	return s.GetConf(s.PlatName, "mqcserver", "mqc", s.ClusterName)
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

//GetRouter 获取服务器的路由配置
func (s *SConf) getRouter(tp string) *services.ORouter {
	switch tp {
	case global.API:
		return s.Service.API
	case global.Web:
		return s.Service.Web
	case global.WS:
		return s.Service.WS
	case global.RPC:
		return s.Service.RPC
	default:
		panic(fmt.Sprintf("无法获取服务%s的路由配置", tp))
	}
}
