package mqc

import (
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

//Middlewares 中间件
var Middlewares = make(middleware.Handlers, 0, 1)
