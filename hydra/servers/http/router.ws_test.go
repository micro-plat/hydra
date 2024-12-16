package http

import (
	x "net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/lib4go/assert"
)

func TestServer_addWSRouters(t *testing.T) {

	tests := []struct {
		name    string
		routers []*router.Router
	}{
		{name: "1. httpserver-添加ws空路由", routers: []*router.Router{}},
		{name: "2. httpserver-添加ws单条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"GET"})}},
		{name: "3. httpserver-添加ws多条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"GET"}), router.NewRouter("/rpcserver/taosy/test1", "/rpcserver/taosy/test1", []string{"POST"})}},
	}

	for _, tt := range tests {
		s := &Server{option: &option{}, server: &x.Server{}}
		opt := WithServerType("ws")
		opt(s.option)
		s.addWSRouters(tt.routers...)
		// assert.Equalf(t, 6, len(s.engine.GetHandlers()), tt.name+",中间件数量")
		assert.Equalf(t, 6, len(s.engine.Routes()), tt.name+",路由数量")
	}
}
