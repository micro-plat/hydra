package middleware

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {

		//检查是否被限流
		ctx.Service(service) //保存服务信息
		if ctx.Request().Path().IsLimited() {
			//降级处理
			fallback(ctx, service)
			return
		}
		//处理RPC服务调用
		if addr, ok := global.IsProto(service, global.ProtoRPC); ok {
			response, err := components.Def.RPC().GetRegularRPC().Swap(addr, ctx)
			if err != nil {
				ctx.Response().Write(response.Status, err)
				return
			}
			ctx.Response().Write(response.Status, response.Result)
			return
		}

		result := services.Def.Call(ctx, service)
		ctx.Response().WriteAny(result)
	}
}
