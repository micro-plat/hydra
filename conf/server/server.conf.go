package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/fsa"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/registry"
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
	GetFSAConf() *fsa.FixedSecretAuth
	GetRASConf() *ras.RASAuth
	GetBasicConf() *basic.BasicAuth
	GetRenderConf() *render.Render
}

var _ IServerConf = &ServerConf{}

//ServerConf 服务器配置信息
type ServerConf struct {
	mainConf conf.IMainConf
	varConf  conf.IVarConf
	header   header.Headers
	jwt      *jwt.JWTAuth
	metric   *metric.Metric
	static   *static.Static
	router   *router.Routers
	fsa      *fsa.FixedSecretAuth
	ras      *ras.RASAuth
	basic    *basic.BasicAuth
	render   *render.Render
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
	s.header = header.GetConf(s.mainConf)
	s.jwt = jwt.GetConf(s.mainConf)
	s.metric = metric.GetConf(s.mainConf)
	s.static = static.GetConf(s.mainConf)
	s.router = router.GetConf(s.mainConf)
	s.fsa = fsa.GetConf(s.mainConf)
	s.ras = ras.GetConf(s.mainConf)
	s.render = render.GetConf(s.mainConf)
	s.basic = basic.GetConf(s.mainConf)
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
	return s.header
}

//GetJWTConf 获取jwt配置
func (s *ServerConf) GetJWTConf() *jwt.JWTAuth {
	return s.jwt
}

//GetMetricConf 获取metric配置
func (s *ServerConf) GetMetricConf() *metric.Metric {
	return s.metric
}

//GetStaticConf 获取静态文件配置
func (s *ServerConf) GetStaticConf() *static.Static {
	return s.static
}

//GetRouterConf 获取路由信息
func (s *ServerConf) GetRouterConf() *router.Routers {
	return s.router
}

//GetFSAConf 获取路由信息
func (s *ServerConf) GetFSAConf() *fsa.FixedSecretAuth {
	return s.fsa
}

//GetRASConf 获取路由信息
func (s *ServerConf) GetRASConf() *ras.RASAuth {
	return s.ras
}

//GetBasicConf 获取basic认证配置
func (s *ServerConf) GetBasicConf() *basic.BasicAuth {
	return s.basic
}

//GetRenderConf 获取状态渲染控件
func (s *ServerConf) GetRenderConf() *render.Render {
	return s.render
}
