package adapter

import (
	"net/http"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

//IEngine IEngine
type IEngine interface {
	Use(handlers ...middleware.Handler)
	Handle(method string, path string, handler middleware.Handler)
	Routes() RoutesInfo
	SetHandlers(handlers ...middleware.Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
	HandleRequest(request IRequest) (response IResponseWriter, err error)
}

//Engine AdapterEngine
type Engine struct {
	handlers      middleware.Handlers
	nodes         []node
	wrapperEngine IEngine
	onceLock      sync.Once
}

//New New
func New(wrapperEngine IEngine) *Engine {
	return &Engine{
		wrapperEngine: wrapperEngine,
		handlers:      []middleware.Handler{},
		nodes:         []node{},
	}
}

//Use Use
func (engine *Engine) Use(handlers ...middleware.Handler) {
	engine.wrapperEngine.Use(handlers...)
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

//Handle Handle
func (engine *Engine) Handle(routers ...IRouter) {
	// nodes := engine.interanlHandle(routers...)
	// for _, n := range nodes {
	// 	engine.wrapperEngine.SetHandlers(n.handlers...)
	// 	for j := range n.actions {
	// 		engine.wrapperEngine.Handle(strings.ToUpper(n.actions[j]), n.path)
	// 	}
	// }
	engine.HandleCustom(middleware.ExecuteHandler(), routers...)
	return
}

//HandleCustom Handle
func (engine *Engine) HandleCustom(handler middleware.Handler, routers ...IRouter) {
	nodes := engine.interanlHandle(routers...)
	for _, n := range nodes {
		engine.wrapperEngine.SetHandlers(n.handlers...)
		for j := range n.actions {
			engine.wrapperEngine.Handle(strings.ToUpper(n.actions[j]), n.path, handler)
		}
	}
	return
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.wrapperEngine.ServeHTTP(w, req)
}

//Find 查找服务
func (engine *Engine) Find(service string) bool {
	return false
}

//HandleRequest 查找服务
func (engine *Engine) HandleRequest(request IRequest) (response IResponseWriter, err error) {
	response, err = engine.wrapperEngine.HandleRequest(request)
	return
}

//GetHandlers GetHandlers
func (engine *Engine) GetHandlers() middleware.Handlers {
	return engine.handlers
}

//Routes Routes
func (engine *Engine) Routes() RoutesInfo {
	return engine.wrapperEngine.Routes()
}

// //DispHandle Handle
// func (engine *Engine) DispHandle(tp string, routers ...IRouter) {
// 	nodes := engine.interanlHandle(routers...)
// 	engine.onceLock.Do(func() {
// 		engine.dispEngine.Use(engine.handlers.DispFunc(tp)...)
// 	})

// 	for _, n := range nodes {
// 		engine.dispEngine.Handlers = n.GetDispHandlers(tp)
// 		for j := range n.actions {
// 			engine.dispEngine.Handle(n.actions[j], n.path, middleware.ExecuteHandler().DispFunc(tp))
// 		}
// 	}
// 	return
// }

// //GinEngine GinEngine
// func (engine *Engine) GinEngine() *gin.Engine {
// 	if engine.ginEngine == nil {
// 		engine.ginEngine = gin.New()
// 	}
// 	return engine.ginEngine
// }

// //DispEngine DispEngine
// func (engine *Engine) DispEngine() *dispatcher.Engine {
// 	if engine.dispEngine == nil {
// 		engine.dispEngine = dispatcher.New()
// 	}
// 	return engine.dispEngine
// }

// func (engine *Engine) Find(path string) bool {
// 	return
// }

// func (engine *Engine) HandleRequest(r IRequest) (w ResponseWriter, err error) {
// 	return
// }
