package wss

import (
	"net/http"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func newDispatcher(serverType string, routers []*router.Router) *adapter.DispatcherEngine {
	engine := adapter.NewDispatcherEngine(serverType)
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logging())
	engine.Use(middleware.Trace())
	engine.Use(middleware.Limit())
	engine.Use(middleware.Delay())
	engine.Use(middleware.APIKeyAuth())
	engine.Use(middleware.RASAuth())
	engine.Use(middleware.JwtAuth())
	engine.Use(middleware.Render())
	engine.Use(middleware.JwtWriter())
	engine.Use(middleware.NewMetric().Handle())
	engine.Handles(routers, middleware.ExecuteHandler())
	return engine
}

func dispatch(engine *adapter.DispatcherEngine, method string, path string, query string, body []byte, headers map[string]string) (int, http.Header, []byte, error) {
	writer, err := engine.HandleRequest(newDispatchRequest(method, path, query, body, headers))
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	return writer.Status(), writer.Header(), writer.Data(), nil
}

func dispatchStream(engine *adapter.DispatcherEngine, method string, path string, query string, body []byte, headers map[string]string, writer dispatcher.ResponseWriter) (int, http.Header, error) {
	w, err := engine.HandleRequestWithWriter(newDispatchRequest(method, path, query, body, headers), writer)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return w.Status(), w.Header(), nil
}
