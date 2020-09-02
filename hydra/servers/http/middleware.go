package http

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var httpmiddlewares []middleware.Handler

func init() {
	httpmiddlewares = make([]middleware.Handler, 0)
}

//Use 添加对Http处理的中间件
func Use(handler middleware.Handler) {
	httpmiddlewares = append(httpmiddlewares, handler)
}
