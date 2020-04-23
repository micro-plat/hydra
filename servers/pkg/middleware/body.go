package middleware

//Body 处理请求的body参数
func Body() Handler {
	return func(ctx IRequest) {
		if body, ok := ctx.GetBody(); ok {
			ctx.Set("__body_", body)
		}
		ctx.Next()
	}
}
