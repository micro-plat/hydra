package limiter

//Option 配置选项
type Option func(*Rule)

//WithAction 设置请求类型
func WithAction(action ...string) Option {
	return func(a *Rule) {
		a.Action = action
	}
}

//WithMaxWait 设置请求允许等待的时长
func WithMaxWait(second int) Option {
	return func(a *Rule) {
		a.MaxWait = second
	}
}

//WithFallback 启用服务降级处理
func WithFallback() Option {
	return func(a *Rule) {
		a.Fallback = true
	}
}

//WithReponse 设置响应内容
func WithReponse(status int, content string) Option {
	return func(a *Rule) {
		a.Resp = &Resp{Status: status, Content: content}
	}
}
