package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/registry"
)

//IServerConf 服务器配置信息
type IServerConf interface {
	GetMainConf() conf.IMainConf
	GetVarConf() conf.IVarConf
	GetJWTConf() (*jwt.JWTAuth, error)
	GetHeaderConf() (header.Headers, error)
	GetMetricConf() (*metric.Metric, error)
	GetStaticConf() (*static.Static, error)
	GetRouterConf() (*router.Routers, error)
	GetCRONTaskConf() (*task.Tasks, error)

	GetMQCMainConf() (*mqc.Server, error)
	GetMQCQueueConf() (*queue.Queues, error)
	GetAPIKeyConf() (*apikey.APIKeyAuth, error)
	GetRASConf() (*ras.RASAuth, error)
	GetBasicConf() (*basic.BasicAuth, error)
	GetRenderConf() (*render.Render, error)
	GetWhiteListConf() (*whitelist.WhiteList, error)
	GetBlackListConf() (*blacklist.BlackList, error)
	GetLimiter() (*limiter.Limiter, error)
	GetGray() (*gray.Gray, error)
	//获取远程日志配置
	GetRLogConf() (*rlog.Layout, error)
	Close() error
}

var _ IServerConf = &ServerConf{}

//ServerConf 服务器配置信息
type ServerConf struct {
	mainConf conf.IMainConf
	varConf  conf.IVarConf
	*httpSub
	*cronSub
	*mqcSub
	*varSub
}

//NewServerConfBy 构建服务器配置缓存
func NewServerConfBy(platName, sysName, serverType, clusterName string, rgst registry.IRegistry) (s *ServerConf, err error) {
	s = &ServerConf{}
	s.mainConf, err = NewMainConf(platName, sysName, serverType, clusterName, rgst)
	if err != nil {
		return nil, err
	}
	s.varConf, err = NewVarConf(s.mainConf.GetVarPath(), rgst)
	if err != nil {
		return nil, err
	}
	s.httpSub = newhttpSub(s.mainConf)
	s.cronSub = newCronSub(s.mainConf)
	s.mqcSub = newMQCSub(s.mainConf)
	s.varSub = newVarSub(s.varConf)
	return s, nil

}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, rgst registry.IRegistry) (s *ServerConf, err error) {
	platName, sysName, serverType, clusterName := Split(mainConfpath)
	return NewServerConfBy(platName, sysName, serverType, clusterName, rgst)

}

//Close 关闭清理资源
func (s *ServerConf) Close() error {
	return s.mainConf.Close()
}

//GetMainConf 获取服务器主配置
func (s *ServerConf) GetMainConf() conf.IMainConf {
	return s.mainConf
}

//GetVarConf 获取变量配置
func (s *ServerConf) GetVarConf() conf.IVarConf {
	return s.varConf
}
