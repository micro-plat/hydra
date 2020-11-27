package mqc

import (
	"testing"

	"github.com/micro-plat/hydra/test/assert"

	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

func TestUse(t *testing.T) {
	tests := []struct {
		name    string
		handler middleware.Handler
		wantLen int
	}{
		{name: "1.添加一个中间件", handler: func(middleware.IMiddleContext) {}, wantLen: 1},
		{name: "2.再次添加一个中间件", handler: func(middleware.IMiddleContext) {}, wantLen: 2},
	}
	for _, tt := range tests {
		Use(tt.handler)
		assert.Equalf(t, tt.wantLen, len(mqcmiddlewares), tt.name)
	}
}
