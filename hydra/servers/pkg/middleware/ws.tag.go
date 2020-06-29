package middleware


//WSTag 添加ws请求标签
func WSTag() Handler {
	return func(ctx IMiddleContext){		
			ctx.Next()
			ctx.Response().AddSpecial("ws")
	}
}