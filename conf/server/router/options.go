package router

//Option 路由配置选项
type Option func(*Router)

//WithPages 设置当前服务对应的页面信息
func WithPages(p ...string) Option {
	return func(a *Router) {
		a.Pages = p
	}
}

//WithEncoding 设置当前服务对应的编码方式
func WithEncoding(encoding string) Option {
	return func(a *Router) {
		a.Encoding = encoding
	}
}
