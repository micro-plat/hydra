package middleware

import (
	"fmt"

	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/errs"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {
		h, ok := services.Registry.GetHandler(ctx.ServerConf().GetMainConf().GetServerType(), service, ctx.Request().Path().GetMethod())
		if !ok {
			ctx.Response().Write(404, fmt.Sprintf("未找到服务：%s", service))
			return
		}
		result := h.Handle(ctx)
		if err := errs.GetError(result); err != nil {
			ctx.Log().Error("err:", err)
		}
		ctx.Response().WriteAny(result)
	}
}
