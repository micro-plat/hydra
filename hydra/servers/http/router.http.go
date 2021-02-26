package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func (s *Server) addHttpRouters(routers ...*router.Router) {
	if !s.ginTrace {
		gin.SetMode(gin.ReleaseMode)
	}
	s.engine = gin.New()
	s.engine.Use(middleware.Recovery().GinFunc(s.serverType))
	s.engine.Use(middleware.Logging().GinFunc())   //记录请求日志
	s.engine.Use(middleware.Processor().GinFunc()) //前缀处理
	s.engine.Use(middleware.Recovery().GinFunc())
	s.engine.Use(s.metric.Handle().GinFunc()) //生成metric报表
	// s.engine.Use(middleware.APM().GinFunc())       //链数跟踪
	s.engine.Use(middleware.Trace().GinFunc())     //跟踪信息
	s.engine.Use(middleware.BlackList().GinFunc()) //黑名单控制
	s.engine.Use(middleware.WhiteList().GinFunc()) //白名单控制
	s.engine.Use(middleware.Proxy().GinFunc())     //灰度配置
	s.engine.Use(middleware.Delay().GinFunc())     //
	s.engine.Use(middleware.Limit().GinFunc())     //限流处理
	s.engine.Use(middleware.Header().GinFunc())    //设置请求头
	s.engine.Use(middleware.Static().GinFunc())    //处理静态文件
	s.engine.Use(middleware.Options().GinFunc())   //处理option响应
	s.engine.Use(middleware.BasicAuth().GinFunc()) //
	s.engine.Use(middleware.APIKeyAuth().GinFunc())
	s.engine.Use(middleware.RASAuth().GinFunc())
	s.engine.Use(middleware.JwtAuth().GinFunc()) //jwt安全认证
	s.engine.Use(middlewares.GinFunc()...)

	s.engine.Use(middleware.Render().GinFunc())    //响应渲染组件
	s.engine.Use(middleware.JwtWriter().GinFunc()) //设置jwt回写

	s.addRouter(routers...)
	s.server.Handler = s.engine
	return
}

func (s *Server) addRouter(routers ...*router.Router) {
	for _, router := range routers {
		for _, method := range router.Action {
			s.engine.Handle(strings.ToUpper(method), router.RealPath, middleware.ExecuteHandler(router.Service).GinFunc())
		}
	}
}
