package middleware

import (
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/lib4go/errs"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {
		h := application.Current().GetHandler(ctx.Server().GetMainConf().GetServerType(), service)
		result := h.Handle(ctx)
		if ctx.Response().Written() {
			return
		}
		if err := errs.GetError(result); err != nil {
			ctx.Response().Write(err.GetCode(), err.GetError().Error())
			return
		}
		ctx.Response().WriteAny(result)
	}
}
