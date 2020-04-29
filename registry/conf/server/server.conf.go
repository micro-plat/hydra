package server

import (
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/registry/conf/server/header"
	"github.com/micro-plat/hydra/registry/conf/server/metric"
	"github.com/micro-plat/hydra/registry/conf/server/router"
	"github.com/micro-plat/hydra/registry/conf/server/static"
)

//IServerConf 服务器配置信息
type IServerConf interface {
	GetMainConf() conf.IMainConf
	GetVarConf() conf.IVarConf
	GetJWTConf() *jwt.JWTAuth
	GetHeaderConf() header.Headers
	GetMetricConf() *metric.Metric
	GetStaticConf() *static.Static
	GetRouterConf() *router.Routers
}

//ServerConf 服务器配置信息
type ServerConf struct {
	mainConf conf.IMainConf
	varConf  conf.IVarConf
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
	return s, nil

}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, rgst registry.IRegistry) (s *ServerConf, err error) {
	platName, sysName, serverType, clusterName := Split(mainConfpath)
	return NewServerConfBy(platName, sysName, serverType, clusterName, rgst)

}

//GetMainConf 获取服务器主配置
func (s *ServerConf) GetMainConf() conf.IMainConf {
	return s.mainConf
}

//GetVarConf 获取变量配置
func (s *ServerConf) GetVarConf() conf.IVarConf {
	return s.varConf
}

//GetHeaderConf 获取响应头配置
func (s *ServerConf) GetHeaderConf() header.Headers {
	return header.GetConf(s.mainConf)
}

//GetJWTConf 获取jwt配置
func (s *ServerConf) GetJWTConf() *jwt.JWTAuth {
	return jwt.GetConf(s.mainConf)
}

//GetMetricConf 获取metric配置
func (s *ServerConf) GetMetricConf() *metric.Metric {
	return metric.GetConf(s.mainConf)
}

//GetStaticConf 获取静态文件配置
func (s *ServerConf) GetStaticConf() *static.Static {
	return static.GetConf(s.mainConf)
}

//GetRouterConf 获取路由信息
func (s *ServerConf) GetRouterConf() *router.Routers {
	return router.GetConf(s.mainConf)
}
