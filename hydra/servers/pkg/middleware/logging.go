package middleware

import (
	"net/http"
	"time"
)

//Logging 记录日志
func Logging() Handler {
	return func(ctx IMiddleContext) {

		//1. 整个服务的开始,记录请求时间与日志
		serverType := ctx.APPConf().GetServerConf().GetServerType()
		start := time.Now()
		path := ctx.Request().Path().GetURL().Path
		ctx.Log().Info(serverType+".request:", ctx.Request().Path().GetMethod(), path, "from", ctx.User().GetClientIP())

		//2. 处理业务
		ctx.Next()

		//3. 将结果刷新到响应流
		ctx.Response().Flush()

		//4. 处理响应日志
		code, _, _ := ctx.Response().GetFinalResponse()
		if code >= http.StatusOK && code < http.StatusBadRequest {
			ctx.Log().Info(serverType+".response:", ctx.Request().Path().GetMethod(), path, code, ctx.Response().GetSpecials(), time.Since(start))
		} else {
			ctx.Log().Error(serverType+".response:", ctx.Request().Path().GetMethod(), path, code, ctx.Response().GetSpecials(), time.Since(start))
		}

		//5.释放资源
		ctx.Close()
	}
}
