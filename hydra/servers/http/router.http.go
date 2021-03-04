package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func (s *Server) addHttpRouters(routers ...*router.Router) {
	if !s.ginTrace {
		gin.SetMode(gin.ReleaseMode)
	}
	s.adapterEngine = adapter.New()
	s.engine = s.adapterEngine.GinEngine()

	//s.engine = gin.New()

	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(middleware.Logging()) //记录请求日志
	//s.adapterEngine.Use(middleware.Processor()) //前缀处理
	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(s.metric.Handle()) //生成metric报表
	// s.adapterEngine.Use(middleware.APM())       //链数跟踪
	s.adapterEngine.Use(middleware.Trace())     //跟踪信息
	s.adapterEngine.Use(middleware.BlackList()) //黑名单控制
	s.adapterEngine.Use(middleware.WhiteList()) //白名单控制
	s.adapterEngine.Use(middleware.Proxy())     //灰度配置
	s.adapterEngine.Use(middleware.Delay())     //
	s.adapterEngine.Use(middleware.Limit())     //限流处理
	s.adapterEngine.Use(middleware.Header())    //设置请求头
	s.adapterEngine.Use(middleware.Static())    //处理静态文件
	s.adapterEngine.Use(middleware.Options())   //处理option响应
	s.adapterEngine.Use(middleware.BasicAuth()) //
	s.adapterEngine.Use(middleware.APIKeyAuth())
	s.adapterEngine.Use(middleware.RASAuth())
	s.adapterEngine.Use(middleware.JwtAuth()) //jwt安全认证
	s.adapterEngine.Use(middlewares...)

	s.adapterEngine.Use(middleware.Render())    //响应渲染组件
	s.adapterEngine.Use(middleware.JwtWriter()) //设置jwt回写

	s.server.Handler = s.engine

	s.addRouter(routers...)
	return
}

func (s *Server) addRouter(routers ...*router.Router) {
	adapterRouters := make([]adapter.IRouter, len(routers))
	for i := range routers {
		adapterRouters[i] = routers[i]
	}
	s.adapterEngine.GinHandle(s.serverType, adapterRouters...)

}
