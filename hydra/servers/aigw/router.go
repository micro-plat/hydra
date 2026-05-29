package aigw

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware/aimiddleware"
)

func (s *Server) addAIRouters(routers ...*router.Router) {
	if !s.ginTrace {
		gin.SetMode(gin.ReleaseMode)
	}
	s.engine = adapter.NewGinEngine(s.serverType)
	s.engine.Use(middleware.Recovery(true))
	s.engine.Use(s.metric.Handle())
	s.engine.Use(middleware.Logging())
	s.engine.Use(middleware.Recovery())
	s.engine.Use(middleware.Trace())
	s.engine.Use(middleware.Limit())
	s.engine.Use(middleware.Header())
	s.engine.Use(middleware.Options())
	s.engine.Use(middlewares...)

	s.server.Handler = s.engine
	s.engine.Handles(routers, aimiddleware.ExecuteHandler())
}
