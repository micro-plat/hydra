package middleware

import (
	"time"
)

//Logging 记录日志
func Logging() Handler {
	return func(ctx IMiddleContext) {

		//1. 整个服务的开始,记录请求时间与日志
		start := time.Now()
		path := ctx.Request().Path().GetPath()
		ctx.Log().Info(ctx.ServerConf().GetMainConf().GetServerType()+".request:", ctx.ServerConf().GetMainConf().GetServerName(), ctx.Request().Path().GetMethod(), path, "from", ctx.User().GetClientIP(), ctx.Request().GetTrace())

		//2. 处理业务
		ctx.Next()

		//3. 处理响应
		if ctx.Response().GetStatusCode() >= 200 && ctx.Response().GetStatusCode() < 400 {
			ctx.Log().Info(ctx.ServerConf().GetMainConf().GetServerType()+".response:", ctx.ServerConf().GetMainConf().GetServerName(), ctx.Request().Path().GetMethod(), path, ctx.Response().GetStatusCode(), ctx.Response().GetSpecials(), time.Since(start), ctx.Response().GetTrace())
		} else {
			ctx.Log().Error(ctx.ServerConf().GetMainConf().GetServerType()+".response:", ctx.ServerConf().GetMainConf().GetServerName(), ctx.Request().Path().GetMethod(), path, ctx.Response().GetStatusCode(), ctx.Response().GetSpecials(), time.Since(start), ctx.Response().GetTrace())
		}

		//4.释放资源
		ctx.Close()
	}
}
