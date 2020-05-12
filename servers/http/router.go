package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/registry/conf/server/router"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
)

func (s *Server) addRouters(routers ...*router.Router) {
	// if !application.IsDebug {
	gin.SetMode(gin.ReleaseMode)
	// }
	s.engine = gin.New()
	s.engine.Use(middleware.Recovery().GinFunc(s.serverType))
	s.engine.Use(middleware.Logging().GinFunc()) //记录请求日志
	s.engine.Use(middleware.Trace().GinFunc())
	// s.engine.Use(s.metric.Handle().GinFunc())      //生成metric报表
	s.engine.Use(middleware.Options().GinFunc())   //处理option响应
	s.engine.Use(middleware.Static().GinFunc())    //处理静态文件
	s.engine.Use(middleware.JwtAuth().GinFunc())   //jwt安全认证
	s.engine.Use(middleware.Response().GinFunc())  //处理返回值
	s.engine.Use(middleware.Header().GinFunc())    //设置请求头
	s.engine.Use(middleware.JwtWriter().GinFunc()) //设置jwt回写
	s.addRouter(routers...)
	s.server.Handler = s.engine
	return
}
func (s *Server) addRouter(routers ...*router.Router) {
	for _, router := range routers {
		if router.Disable {
			continue
		}
		for _, method := range router.Action {
			s.engine.Handle(method, router.Path, middleware.ExecuteHandler(router.Service).GinFunc())
		}
	}
}
