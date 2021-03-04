package adapter

import (
	"net/http"

	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var _ IEngine = &EngineWrapperDisp{}

//EngineWrapperDisp EngineWrapperDisp
type EngineWrapperDisp struct {
	serverType string
	engine     *dispatcher.Engine
	handlers   middleware.Handlers
}

//NewEngineWrapperDisp NewEngineWrapperDisp
func NewEngineWrapperDisp(engine *dispatcher.Engine, serverType string) IEngine {
	return &EngineWrapperDisp{
		serverType: serverType,
		engine:     engine,
		handlers:   make(middleware.Handlers, 0),
	}
}

//Use Use
func (e *EngineWrapperDisp) Use(handlers ...middleware.Handler) {
	e.handlers = append(e.handlers, handlers...)
	handlerList := middleware.Handlers(handlers)
	e.engine.Use(handlerList.DispFunc(e.serverType)...)
}

//Handle Handle
func (e *EngineWrapperDisp) Handle(method string, path string, handler middleware.Handler) {
	e.engine.Handle(method, path, handler.DispFunc(e.serverType))

}

//SetHandlers SetHandlers
func (e *EngineWrapperDisp) SetHandlers(handlers ...middleware.Handler) {
	handlerList := middleware.Handlers(handlers)
	e.engine.Handlers = handlerList.DispFunc(e.serverType)
}

//Routes Routes
func (e *EngineWrapperDisp) Routes() RoutesInfo {
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

func (e *EngineWrapperDisp) ServeHTTP(http.ResponseWriter, *http.Request) {

}

//HandleRequest 处理请求
func (e *EngineWrapperDisp) HandleRequest(request IRequest) (response IResponseWriter, err error) {
	response, err = e.engine.HandleRequest(request)
	return
}
