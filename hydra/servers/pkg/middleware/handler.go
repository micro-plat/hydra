package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

// ExecuteHandler 业务处理Handler
func ExecuteHandler() Handler {
	return func(ctx IMiddleContext) {
		service := ctx.Request().Path().GetService()
		//检查是否被限流
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
			headers := response.GetHeaders()
			for k := range headers {
				ctx.Response().Header(k, headers.GetString(k))
			}
			ctx.Response().Write(response.GetStatus(), response.GetResult()) //写入响应内容
			return
		}

		serverType := ctx.APPConf().GetServerConf().GetServerType()
		method := ctx.Request().Path().GetMethod()

		if services.Def.Has(serverType, service, method) {
			result := services.Def.Call(ctx, service)
			if ok, r := context.IsSSEData(result); ok { //处理SSE协议数据
				// ctx.Response().NoNeedWrite(200)
				r.LoopWrite(ctx.Response().GetHTTPReponse())
				return
			}
			ctx.Response().WriteAny(result) //写入响应内容
			return
		}

		ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("未找到路径:%s", ctx.Request().Path().GetRequestPath()))
	}
}
