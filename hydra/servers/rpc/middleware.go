package rpc

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var rpcmiddlewares []middleware.Handler

func init() {
	rpcmiddlewares = make([]middleware.Handler, 0)
}

//Use 添加对Cron处理的中间件
func Use(handler middleware.Handler) {
	rpcmiddlewares = append(rpcmiddlewares, handler)
}
