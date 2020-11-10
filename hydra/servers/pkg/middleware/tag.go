package middleware

//Tag 添加服务标签
func Tag() Handler {
	return func(ctx IMiddleContext) {
		ctx.Next()
		ctx.Response().AddSpecial(ctx.APPConf().GetServerConf().GetServerType())
	}
}
