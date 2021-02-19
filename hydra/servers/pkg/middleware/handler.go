package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

//ExecuteHandler 业务处理Handler
func ExecuteHandler(service string) Handler {
	return func(ctx IMiddleContext) {
		ctx.Log().Info("services.ExecuteHandler")
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
				ctx.Response().Write(response.GetStatus(), err)
				return
			}
			ctx.Response().Write(response.GetStatus(), response.GetResult())
			return
		}

		//处理option
		if checkOption(ctx) {
			return
		}

		//检查服务中是否包含当前请求路径
		if services.Def.Has(ctx.APPConf().GetServerConf().GetServerType(), service) {
			ctx.Log().Debug("services.Def.Call")
			result := services.Def.Call(ctx, service)
			ctx.Response().WriteAny(result)
			return
		}

		//处理静态文件
		ctx.Log().Debug("doStatic")
		if doStatic(ctx, service) {
			return
		}
		ctx.Log().Debug("404")
		ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("未找到路径:%s", ctx.Request().Path().GetRequestPath()))
	}
}
