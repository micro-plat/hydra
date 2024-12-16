package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/http/ws"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func (s *Server) addWSRouters(routers ...*router.Router) {
	if !s.ginTrace {
		gin.SetMode(gin.ReleaseMode)
	}
	s.engine = adapter.NewGinEngine(s.serverType)
	s.engine.Use(middleware.Recovery(true))
	s.engine.Use(middleware.Logging()) //记录请求日志
	s.engine.Use(middleware.Recovery())
	s.engine.Use(middleware.BlackList()) //黑名单控制
	s.engine.Use(middleware.WhiteList()) //白名单控制
	s.engine.Use(middleware.Limit())     //限流处理
	s.engine.Use()
	s.addWSRouter(routers...)
	s.server.Handler = s.engine
}

func (s *Server) addWSRouter(routers ...*router.Router) {
	ws.InitWSEngine(routers...)
	router := router.GetWSHomeRouter()
	for _, a := range router.Action {
		s.engine.Handle(a, router.Path, ws.WSExecuteHandler())
	}

}
