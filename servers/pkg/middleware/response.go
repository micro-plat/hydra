package middleware

//Response 处理服务器响应
func Response() Handler {
	return func(r IMiddleContext) {
		r.Next()
	}
}
