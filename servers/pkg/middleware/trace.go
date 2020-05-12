package middleware

import "strings"

//Trace 系统跟踪日志
func Trace() Handler {
	return func(ctx IMiddleContext) {

		//1. 打印请求参数
		enable := ctx.ServerConf().GetMainConf().GetMainConf().GetBool("trace")
		if enable && strings.ToLower(ctx.Request().Path().GetMethod()) != "get" {
			ctx.Response().AddSpecial("trace")
			body, _ := ctx.Request().GetBody()
			ctx.Log().Info(body)
		}

		//2. 业务处理
		ctx.Next()

		//3. 打印响应参数
		if enable {
			ctx.Log().Info(ctx.Response().GetResponse())
		}
	}
}
