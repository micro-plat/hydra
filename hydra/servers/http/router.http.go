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
	s.engine = adapter.NewGinEngine(s.serverType)
	s.engine.Use(middleware.Recovery(true))
	s.engine.Use(s.metric.Handle()) //生成metric报表
	s.engine.Use(middleware.Gzip(middleware.DefaultCompression))
	s.engine.Use(middleware.Logging()) //记录请求日志
	s.engine.Use(middleware.Recovery())

	// s.engine.Use(middleware.APM())       //链数跟踪
	s.engine.Use(middleware.Trace())     //跟踪信息
	s.engine.Use(middleware.BlackList()) //黑名单控制
	s.engine.Use(middleware.WhiteList()) //白名单控制
	s.engine.Use(middleware.Proxy())     //灰度配置
	s.engine.Use(middleware.Delay())     //
	s.engine.Use(middleware.Limit())     //限流处理
	s.engine.Use(middleware.Header())    //设置请求头
	s.engine.Use(middleware.Static())    //处理静态文件
	s.engine.Use(middleware.Options())   //处理option响应
	s.engine.Use(middleware.BasicAuth()) //
	s.engine.Use(middleware.APIKeyAuth())
	s.engine.Use(middleware.RASAuth())
	s.engine.Use(middleware.JwtAuth()) //jwt安全认证
	s.engine.Use(middlewares...)

	s.engine.Use(middleware.Render())    //响应渲染组件
	s.engine.Use(middleware.JwtWriter()) //设置jwt回写

	s.server.Handler = s.engine

	s.addRouter(routers...)
}

func (s *Server) addRouter(routers ...*router.Router) {
	s.engine.Handles(routers, middleware.ExecuteHandler())
}
