package middleware

//Trace 系统跟踪日志
func Trace() Handler {
	return func(ctx IMiddleContext) {

		enable := ctx.APPConf().GetServerConf().GetMainConf().GetBool("trace")
		if !enable {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("trace")

		//1.打印请求参数
		input, _ := ctx.Request().GetMap()
		ctx.Log().Debug("> trace.request:", input)

		//2. 业务处理
		ctx.Next()

		//3. 打印响应参数
		s, c := ctx.Response().GetFinalResponse()
		ctx.Log().Debug("> trace.response:", s, c)

	}
}
