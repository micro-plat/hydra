package middleware

import (
	"net/http"
	"time"
)

//Logging 记录日志
func Logging() Handler {
	return func(ctx IMiddleContext) {

		//1. 整个服务的开始,记录请求时间与日志
		start := time.Now()
		path := ctx.Request().Path().GetURL()
		ctx.Log().Info(ctx.APPConf().GetServerConf().GetServerType()+".request:", ctx.Request().Path().GetMethod(), path, "from", ctx.User().GetClientIP())

		//2. 处理业务
		ctx.Next()

		//3. 将结果刷新到响应流
		ctx.Response().Flush()

		//4. 处理响应日志
		code, _ := ctx.Response().GetFinalResponse()
		if code >= http.StatusOK && code < http.StatusBadRequest {
			ctx.Log().Info(ctx.APPConf().GetServerConf().GetServerType()+".response:", ctx.Request().Path().GetMethod(), path, code, ctx.Response().GetSpecials(), time.Since(start))
		} else {
			ctx.Log().Error(ctx.APPConf().GetServerConf().GetServerType()+".response:", ctx.Request().Path().GetMethod(), path, code, ctx.Response().GetSpecials(), time.Since(start))
		}

		//5.释放资源
		ctx.Close()
	}
}
