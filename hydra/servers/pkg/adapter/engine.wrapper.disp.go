package adapter

import (
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

//DispatcherEngine EngineWrapperGin
type DispatcherEngine struct {
	serverType string
	*dispatcher.Engine
}

//NewDispatcherEngine NewEngineWrapperGin
func NewDispatcherEngine(serverType string) *DispatcherEngine {
	return &DispatcherEngine{
		serverType: serverType,
		Engine:     dispatcher.New(),
	}
}

//Use Use
func (e *DispatcherEngine) Use(handlers ...middleware.Handler) {
	for _, h := range handlers {
		e.Engine.Use(h.DispFunc(e.serverType))
	}
}
func (e *DispatcherEngine) Handles(routers []*router.Router, handler middleware.Handler, hds ...middleware.Handler) {
	for _, r := range routers {
		for _, action := range r.Action {
			e.Handle(strings.ToUpper(action), r.Path, middleware.ExecuteHandler())
		}
	}
}

//Handle Handle
func (e *DispatcherEngine) Handle(method string, path string, handler middleware.Handler, hds ...middleware.Handler) {
	handlers := make([]dispatcher.HandlerFunc, 0, len(hds)+1)
	for _, h := range hds {
		handlers = append(handlers, h.DispFunc(e.serverType))
	}
	handlers = append(handlers, handler.DispFunc(e.serverType))
	e.Engine.Handle(method, path, handlers...)
}
