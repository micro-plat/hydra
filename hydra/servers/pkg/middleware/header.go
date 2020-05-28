package middleware

var originName = "Origin"

//Header 响应头设置
func Header() Handler {

	return func(ctx IMiddleContext) {

		//1. 获取header配置
		headers := ctx.ServerConf().GetHeaderConf()
		if len(headers) > 0 {
			ctx.Response().AddSpecial("hdr")
		}

		//3. 处理响应header参数
		origin := ctx.Request().Path().GetHeader(originName)
		hds := headers.GetHeaderByOrigin(origin)
		for k, v := range hds {
			ctx.Response().SetHeader(k, v)
		}

		//2. 业务处理
		ctx.Next()

	}
}
