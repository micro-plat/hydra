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
	s.header = GetLoader(cnf, s.getHeaderConfFunc())
	s.jwt = GetLoader(cnf, s.getJWTConfFunc())
	s.metric = GetLoader(cnf, s.getMetricConfFunc())
	s.static = GetLoader(cnf, s.getStaticConfFunc())
	s.router = GetLoader(cnf, s.getRouterConfFunc())
	s.apikey = GetLoader(cnf, s.getAPIKeyConfFunc())
	s.ras = GetLoader(cnf, s.getRasFunc())
	s.basic = GetLoader(cnf, s.getBasicFunc())
	s.render = GetLoader(cnf, s.getRenderFunc())
	s.whiteList = GetLoader(cnf, s.getWhitelistFunc())
	s.blackList = GetLoader(cnf, s.getBlacklistFunc())
	s.limit = GetLoader(cnf, s.getLimiterFunc())
	s.gray = GetLoader(cnf, s.getGrayFunc())
	return s
}

//getHeaderConfFunc 获取header配置信息
func (s httpSub) getHeaderConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return header.GetConf(cnf)
	}
}

//getJWTConfFunc 获取jwt配置信息
func (s httpSub) getJWTConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return jwt.GetConf(cnf)
	}
}

//getMetricConfFunc 获取metric配置信息
func (s httpSub) getMetricConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return metric.GetConf(cnf)
	}
}

//getStaticConfFunc 获取static配置信息
func (s httpSub) getStaticConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return static.GetConf(cnf)
	}
}

//getRouterConfFunc 获取router配置信息
func (s httpSub) getRouterConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return router.GetConf(cnf)
	}
}

//getAPIKeyConfFunc 获取apikey配置信息
func (s httpSub) getAPIKeyConfFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return apikey.GetConf(cnf)
	}
}

//getRasFunc 获取ras配置信息
func (s httpSub) getRasFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return ras.GetConf(cnf)
	}
}

//getBasicFunc 获取basic配置信息
func (s httpSub) getBasicFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return basic.GetConf(cnf)
	}
}

//getRenderFunc 获取render配置信息
func (s httpSub) getRenderFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return render.GetConf(cnf)
	}
}

//getWhitelistFunc 获取whitelist配置信息
func (s httpSub) getWhitelistFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return whitelist.GetConf(cnf)
	}
}

//getBlacklistFunc 获取blacklist配置信息
func (s httpSub) getBlacklistFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return blacklist.GetConf(cnf)
	}
}

//getLimiterFunc 获取limiter配置信息
func (s httpSub) getLimiterFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return limiter.GetConf(cnf)
	}
}

//getGrayFunc 获取gray配置信息
func (s httpSub) getGrayFunc() func(cnf conf.IMainConf) (interface{}, error) {
	return func(cnf conf.IMainConf) (interface{}, error) {
		return gray.GetConf(cnf)
	}
}

//GetHeaderConf 获取响应头配置
func (s *httpSub) GetHeaderConf() (header.Headers, error) {
	headerObj, err := s.header.GetConf()
	if err != nil {
		return nil, err
	}

	return headerObj.(header.Headers), nil
}

//GetJWTConf 获取jwt配置
func (s *httpSub) GetJWTConf() (*jwt.JWTAuth, error) {
	jwtObj, err := s.jwt.GetConf()
	if err != nil {
		return nil, err
	}
	return jwtObj.(*jwt.JWTAuth), nil
}

//GetMetricConf 获取metric配置
func (s *httpSub) GetMetricConf() (*metric.Metric, error) {
	metricObj, err := s.metric.GetConf()
	if err != nil {
		return nil, err
	}
	return metricObj.(*metric.Metric), nil
}

//GetStaticConf 获取静态文件配置
func (s *httpSub) GetStaticConf() (*static.Static, error) {
	staticObj, err := s.static.GetConf()
	if err != nil {
		return nil, err
	}
	return staticObj.(*static.Static), nil
}

//GetRouterConf 获取路由信息
func (s *httpSub) GetRouterConf() (*router.Routers, error) {
	routerObj, err := s.router.GetConf()
	if err != nil {
		return nil, err
	}

	return routerObj.(*router.Routers), nil
}

//GetAPIKeyConf 获取apikey配置
func (s *httpSub) GetAPIKeyConf() (*apikey.APIKeyAuth, error) {
	apikeyObj, err := s.apikey.GetConf()
	if err != nil {
		return nil, err
	}

	return apikeyObj.(*apikey.APIKeyAuth), nil
}

//GetRASConf 获取RAS配置信息
func (s *httpSub) GetRASConf() (*ras.RASAuth, error) {
	rasObj, err := s.ras.GetConf()
	if err != nil {
		return nil, err
	}

	return rasObj.(*ras.RASAuth), nil
}

//GetBasicConf 获取basic认证配置
func (s *httpSub) GetBasicConf() (*basic.BasicAuth, error) {
	basicObj, err := s.basic.GetConf()
	if err != nil {
		return nil, err
	}
	return basicObj.(*basic.BasicAuth), nil
}

//GetRenderConf 获取状态渲染控件
func (s *httpSub) GetRenderConf() (*render.Render, error) {
	renderObj, err := s.render.GetConf()
	if err != nil {
		return nil, err
	}
	return renderObj.(*render.Render), nil
}

//GetWhiteListConf 获取白名单配置
func (s *httpSub) GetWhiteListConf() (*whitelist.WhiteList, error) {
	whiteListObj, err := s.whiteList.GetConf()
	if err != nil {
		return nil, err
	}

	return whiteListObj.(*whitelist.WhiteList), nil
}

//GetBlackListConf 获取黑名单配置
func (s *httpSub) GetBlackListConf() (*blacklist.BlackList, error) {
	blackListObj, err := s.blackList.GetConf()
	if err != nil {
		return nil, err
	}
	return blackListObj.(*blacklist.BlackList), nil
}

//GetLimiter 获取限流配置
func (s *httpSub) GetLimiter() (*limiter.Limiter, error) {
	limitObj, err := s.limit.GetConf()
	if err != nil {
		return nil, err
	}
	return limitObj.(*limiter.Limiter), nil
}

//GetGray 获取灰度配置
func (s *httpSub) GetGray() (*gray.Gray, error) {
	grayObj, err := s.gray.GetConf()
	if err != nil {
		return nil, err
	}

	return grayObj.(*gray.Gray), nil
}
