package adapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var _ IEngine = &EngineWrapperGin{}

//EngineWrapperGin EngineWrapperGin
type EngineWrapperGin struct {
	serverType string
	engine     *gin.Engine
	handlers   middleware.Handlers
}

//NewEngineWrapperGin NewEngineWrapperGin
func NewEngineWrapperGin(engine *gin.Engine, serverType string) IEngine {
	return &EngineWrapperGin{
		serverType: serverType,
		engine:     engine,
		handlers:   make(middleware.Handlers, 0),
	}
}

//Use Use
func (e *EngineWrapperGin) Use(handlers ...middleware.Handler) {
	e.handlers = append(e.handlers, handlers...)
	handlerList := middleware.Handlers(handlers)
	e.engine.Use(handlerList.GinFunc(e.serverType)...)
}

//Handle Handle
func (e *EngineWrapperGin) Handle(method string, path string, handler middleware.Handler) {
	e.engine.Handle(method, path, handler.GinFunc(e.serverType))
}

//SetHandlers SetHandlers
func (e *EngineWrapperGin) SetHandlers(handlers ...middleware.Handler) {
	handlerList := middleware.Handlers(handlers)
	e.engine.Handlers = handlerList.GinFunc(e.serverType)
}

//Routes Routes
func (e *EngineWrapperGin) Routes() RoutesInfo {
	erts := e.engine.Routes()
	result := make(RoutesInfo, len(erts))

	for i := range erts {
		cur := erts[i]
		result[i] = RouteInfo{
			Method:  cur.Method,
			Path:    cur.Path,
			Handler: cur.Handler,
		}
	}
	return result
}

func (e *EngineWrapperGin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.engine.ServeHTTP(w, req)
}

//HandleRequest 查找服务
func (e *EngineWrapperGin) HandleRequest(request IRequest) (response IResponseWriter, err error) {

	return
}
