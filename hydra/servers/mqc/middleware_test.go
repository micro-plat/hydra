package mqc

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func TestUse(t *testing.T) {
	tests := []struct {
		name    string
		handler []middleware.Handler
		wantLen int
	}{
		{name: "1. mqc服务中间件添加-空列表", handler: []middleware.Handler{}, wantLen: 0},
		{name: "2. mqc服务中间件添加-一个中间件", handler: []middleware.Handler{func(middleware.IMiddleContext) {}}, wantLen: 1},
		{name: "3. mqc服务中间件添加-多个中间件", handler: []middleware.Handler{func(middleware.IMiddleContext) {}, func(middleware.IMiddleContext) {}}, wantLen: 2},
	}
	for _, tt := range tests {
		mqcmiddlewares = make([]middleware.Handler, 0)
		for _, h := range tt.handler {
			Use(h)
		}
		assert.Equalf(t, tt.wantLen, len(mqcmiddlewares), tt.name)
	}
}
