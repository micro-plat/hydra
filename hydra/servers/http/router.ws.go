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
	s.adapterEngine = adapter.New(adapter.NewEngineWrapperGin(gin.New(), s.serverType))

	//s.engine = gin.New()
	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(middleware.Logging()) //记录请求日志
	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(middleware.BlackList()) //黑名单控制
	s.adapterEngine.Use(middleware.WhiteList()) //白名单控制
	s.adapterEngine.Use(middleware.Limit())     //限流处理
	s.adapterEngine.Use()
	s.addWSRouter(routers...)
	s.server.Handler = s.adapterEngine
	return
}

func (s *Server) addWSRouter(routers ...*router.Router) {
	ws.InitWSEngine(routers...)
	router := router.GetWSHomeRouter()
	s.adapterEngine.HandleCustom(ws.WSExecuteHandler(), router)
}
