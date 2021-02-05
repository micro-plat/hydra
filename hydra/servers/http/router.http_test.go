package http

import (
	x "net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/lib4go/assert"
)

func TestServer_addHttpRouters(t *testing.T) {
	tests := []struct {
		name    string
		routers []*router.Router
	}{
		{name: "1. httpserver-添加http空路由", routers: []*router.Router{}},
		{name: "2. httpserver-添加http单条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"Get"})}},
		{name: "3. httpserver-添加http多条路由", routers: []*router.Router{router.NewRouter("/rpcserver/taosy/test", "/rpcserver/taosy/test", []string{"Get"}), router.NewRouter("/rpcserver/taosy/test1", "/rpcserver/taosy/test1", []string{"Post"})}},
	}

	for _, tt := range tests {
		s := &Server{option: &option{}, server: &x.Server{}}
		opt := WithServerType("api")
		opt(s.option)
		s.addHttpRouters(tt.routers...)
		assert.Equalf(t, 19, len(s.engine.RouterGroup.Handlers), tt.name+",中间件数量")
		assert.Equalf(t, len(tt.routers), len(s.engine.Routes()), tt.name+",路由数量")
	}
}
