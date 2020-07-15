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
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/static"
)

type httpSub struct {
	cnf       conf.IMainConf
	header    *Loader
	jwt       *Loader
	metric    *Loader
	static    *Loader
	router    *Loader
	apikey    *Loader
	ras       *Loader
	basic     *Loader
	render    *Loader
	whiteList *Loader
	blackList *Loader
	limit     *Loader
	gray      *Loader
}

func newhttpSub(cnf conf.IMainConf) *httpSub {
	s := &httpSub{cnf: cnf}
	s.header = GetLoader(cnf, header.ConfHandler(header.GetConf).Handle)
	s.jwt = GetLoader(cnf, jwt.ConfHandler(jwt.GetConf).Handle)
	s.metric = GetLoader(cnf, metric.ConfHandler(metric.GetConf).Handle)
	s.static = GetLoader(cnf, static.ConfHandler(static.GetConf).Handle)
	s.router = GetLoader(cnf, router.ConfHandler(router.GetConf).Handle)
	s.apikey = GetLoader(cnf, apikey.ConfHandler(apikey.GetConf).Handle)
	s.ras = GetLoader(cnf, ras.ConfHandler(ras.GetConf).Handle)
	s.basic = GetLoader(cnf, basic.ConfHandler(basic.GetConf).Handle)
	s.render = GetLoader(cnf, render.ConfHandler(render.GetConf).Handle)
	s.whiteList = GetLoader(cnf, whitelist.ConfHandler(whitelist.GetConf).Handle)
	s.blackList = GetLoader(cnf, blacklist.ConfHandler(blacklist.GetConf).Handle)
	s.limit = GetLoader(cnf, limiter.ConfHandler(limiter.GetConf).Handle)
	s.gray = GetLoader(cnf, gray.ConfHandler(gray.GetConf).Handle)
	return s
}

//GetHeaderConf 获取响应头配置
func (s *httpSub) GetHeaderConf() header.Headers {
	return s.header.GetConf().(header.Headers)
}

//GetJWTConf 获取jwt配置
func (s *httpSub) GetJWTConf() *jwt.JWTAuth {
	return s.jwt.GetConf().(*jwt.JWTAuth)
}

//GetMetricConf 获取metric配置
func (s *httpSub) GetMetricConf() *metric.Metric {
	return s.metric.GetConf().(*metric.Metric)
}

//GetStaticConf 获取静态文件配置
func (s *httpSub) GetStaticConf() *static.Static {
	return s.static.GetConf().(*static.Static)
}

//GetRouterConf 获取路由信息
func (s *httpSub) GetRouterConf() *router.Routers {
	return s.router.GetConf().(*router.Routers)
}

//GetAPIKeyConf 获取apikey配置
func (s *httpSub) GetAPIKeyConf() *apikey.APIKeyAuth {
	return s.apikey.GetConf().(*apikey.APIKeyAuth)
}

//GetRASConf 获取RAS配置信息
func (s *httpSub) GetRASConf() *ras.RASAuth {
	return s.ras.GetConf().(*ras.RASAuth)
}

//GetBasicConf 获取basic认证配置
func (s *httpSub) GetBasicConf() *basic.BasicAuth {
	return s.basic.GetConf().(*basic.BasicAuth)
}

//GetRenderConf 获取状态渲染控件
func (s *httpSub) GetRenderConf() *render.Render {
	return s.render.GetConf().(*render.Render)
}

//GetWhiteListConf 获取白名单配置
func (s *httpSub) GetWhiteListConf() *whitelist.WhiteList {
	return s.whiteList.GetConf().(*whitelist.WhiteList)
}

//GetBlackListConf 获取黑名单配置
func (s *httpSub) GetBlackListConf() *blacklist.BlackList {
	return s.blackList.GetConf().(*blacklist.BlackList)
}

//GetLimiter 获取限流配置
func (s *httpSub) GetLimiter() *limiter.Limiter {
	return s.limit.GetConf().(*limiter.Limiter)
}

//GetGray 获取灰度配置
func (s *httpSub) GetGray() *gray.Gray {
	return s.gray.GetConf().(*gray.Gray)
}
