package middleware

//ServerType 处理服务器类型，后端通过服务器类型获取服务器配置
func ServerType(tp string) Handler {
	return func(ctx IMiddleContext) {
		ctx.Server().SetServerType(tp)
		ctx.Next()
	}
}
