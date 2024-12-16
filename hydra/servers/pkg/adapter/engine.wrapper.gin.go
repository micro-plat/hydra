package adapter

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

//GinEngine GinEngine
type GinEngine struct {
	serverType string
	*gin.Engine
}

//NewGinEngine NewGinEngine
func NewGinEngine(serverType string) *GinEngine {
	return &GinEngine{
		serverType: serverType,
		Engine:     gin.New(),
	}
}

//Use Use
func (e *GinEngine) Use(handlers ...middleware.Handler) {
	for _, h := range handlers {
		e.Engine.Use(h.GinFunc(e.serverType))
	}
}

//Handle Handle
func (e *GinEngine) Handle(method string, path string, handler middleware.Handler, hds ...middleware.Handler) {
	handlers := make([]gin.HandlerFunc, 0, len(hds)+1)
	for _, h := range hds {
		handlers = append(handlers, h.GinFunc(e.serverType))
	}
	handlers = append(handlers, handler.GinFunc(e.serverType))
	e.Engine.Handle(method, path, handlers...)
}

//Handles Handles
func (e *GinEngine) Handles(routers []*router.Router, handler middleware.Handler, hds ...middleware.Handler) {
	for _, r := range routers {
		for _, action := range r.Action {
			e.Handle(strings.ToUpper(action), r.Path, handler, hds...)
		}
	}
}
