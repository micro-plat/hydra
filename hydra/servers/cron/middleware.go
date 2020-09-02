package cron

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var cronmiddlewares []middleware.Handler

func init() {
	cronmiddlewares = make([]middleware.Handler, 0)
}

//Use 添加对Cron处理的中间件
func Use(handler middleware.Handler) {
	cronmiddlewares = append(cronmiddlewares, handler)
}
