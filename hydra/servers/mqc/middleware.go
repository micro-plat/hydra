package mqc

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var mqcmiddlewares []middleware.Handler

func init() {
	mqcmiddlewares = make([]middleware.Handler, 0)
}

//Use 添加对mqc处理的中间件
func Use(handler middleware.Handler) {
	mqcmiddlewares = append(mqcmiddlewares, handler)
}
