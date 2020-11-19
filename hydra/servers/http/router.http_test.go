package http

import (
	x "net/http"
	"testing"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/test/assert"
)

func TestServer_addHttpRouters(t *testing.T) {
	tests := []struct {
		name    string
		routers []*router.Router
	}{
		{name: "添加http的中间件和路由", routers: []*router.Router{}},
	}

	for _, tt := range tests {
		s := &Server{option: &option{}, server: &x.Server{}}
		opt := WithServerType("api")
		opt(s.option)
		s.addHttpRouters(tt.routers...)
		assert.Equal(t, 18, len(s.engine.RouterGroup.Handlers), tt.name)
	}
}
