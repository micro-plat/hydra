package adapter

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

//Engine AdapterEngine
type Engine struct {
	handlers   middleware.Handlers
	nodes      []node
	ginEngine  *gin.Engine
	dispEngine *dispatcher.Engine
	onceLock   sync.Once
}

//New New
func New() *Engine {
	return &Engine{
		handlers: []middleware.Handler{},
		nodes:    []node{},
	}
}

//Use Use
func (engine *Engine) Use(handlers ...middleware.Handler) {
	engine.handlers = append(engine.handlers, handlers...)
	return
}

//Handle Handle
func (engine *Engine) interanlHandle(routers ...IRouter) []node {
	nodes := make([]node, 0)
	for _, router := range routers {
		curHandlers := make([]middleware.Handler, len(engine.handlers)+1)
		curHandlers[0] = middleware.Service(router.GetService())
		copy(curHandlers[1:], engine.handlers)
		nodes = append(nodes, node{
			path:     router.GetPath(),
			service:  router.GetService(),
			actions:  router.GetActions(),
			handlers: curHandlers,
		})
	}
	return nodes
}

//GinHandle Handle
func (engine *Engine) GinHandle(tp string, routers ...IRouter) {
	nodes := engine.interanlHandle(routers...)
	engine.onceLock.Do(func() {
		engine.ginEngine.Use(engine.handlers.GinFunc(tp)...)
	})
	for _, n := range nodes {
		engine.ginEngine.Handlers = n.GetGinHandlers(tp)
		for j := range n.actions {
			engine.ginEngine.Handle(strings.ToUpper(n.actions[j]), n.path, middleware.ExecuteHandler().GinFunc(tp))
		}
	}
	return
}

//DispHandle Handle
func (engine *Engine) DispHandle(tp string, routers ...IRouter) {
	nodes := engine.interanlHandle(routers...)
	engine.onceLock.Do(func() {
		engine.dispEngine.Use(engine.handlers.DispFunc(tp)...)
	})

	for _, n := range nodes {
		engine.dispEngine.Handlers = n.GetDispHandlers(tp)
		for j := range n.actions {
			engine.dispEngine.Handle(n.actions[j], n.path, middleware.ExecuteHandler().DispFunc(tp))
		}
	}
	return
}

//GinEngine GinEngine
func (engine *Engine) GinEngine() *gin.Engine {
	if engine.ginEngine == nil {
		engine.ginEngine = gin.New()
	}
	return engine.ginEngine
}

//DispEngine DispEngine
func (engine *Engine) DispEngine() *dispatcher.Engine {
	if engine.dispEngine == nil {
		engine.dispEngine = dispatcher.New()
	}
	return engine.dispEngine
}

// func (engine *Engine) Find(path string) bool {
// 	return
// }

// func (engine *Engine) HandleRequest(r IRequest) (w ResponseWriter, err error) {
// 	return
// }
