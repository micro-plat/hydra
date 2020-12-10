package middleware

//APM 跟踪数据
func APM() Handler {
	return func(ctx IMiddleContext) {
		if !ctx.Tracer().Available() {
			ctx.Next()
			return
		}
		ctx.Response().AddSpecial("apm")
		ctx.Tracer().Start()
		defer ctx.Tracer().End()
		ctx.Next()
	}
}
