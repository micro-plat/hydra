package http

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var middlewares = make(middleware.Handlers, 0, 1)

//Middlewares 用户自定义中间件
var Middlewares middleware.ICustomMiddleware = middlewares
