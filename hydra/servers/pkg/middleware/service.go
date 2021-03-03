package middleware

//Service 服务设置处理
func Service(service string) Handler {
	return func(ctx IMiddleContext) {
		ctx.Service(service)
		ctx.Next()
	}
}
