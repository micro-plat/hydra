package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func (s *Server) addWSRouters(routers ...*router.Router) {
	if !global.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.New()
	s.engine.Use(middleware.Recovery().GinFunc(s.serverType))
	s.engine.Use(middleware.Logging().GinFunc())   //记录请求日志
	s.engine.Use(middleware.BlackList().GinFunc()) //黑名单控制
	s.engine.Use(middleware.WhiteList().GinFunc()) //白名单控制
	s.engine.Use(middleware.Limit().GinFunc())     //限流处理
	s.addWSRouter(routers...)
	s.server.Handler = s.engine
	return
}
func (s *Server) addWSRouter(routers ...*router.Router) {

	middleware.InitWSInternalEngine(routers...)

	router := router.GetWSHomeRouter()
	for _, method := range router.Action {
		s.engine.Handle(strings.ToUpper(method), router.Path, middleware.WSExecuteHandler(router.Service).GinFunc())
	}
}
