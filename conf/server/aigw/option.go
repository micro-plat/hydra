package aigw

// Option 配置选项
type Option func(*Server)

// WithTimeout 设置服务器超时时长
func WithTimeout(readTimeout int, writeTimeout int, readHeaderTimeout ...int) Option {
	return func(o *Server) {
		o.RTimeout = readTimeout
		o.WTimeout = writeTimeout
		if len(readHeaderTimeout) > 0 {
			o.RHTimeout = readHeaderTimeout[0]
		}
	}
}

// WithHeaderReadTimeout 设置请求头读取超时时间
func WithHeaderReadTimeout(timeout int) Option {
	return func(o *Server) {
		o.RHTimeout = timeout
	}
}

// WithStreamTimeout 设置流式响应超时时间
func WithStreamTimeout(timeout int) Option {
	return func(o *Server) {
		o.StreamTimeout = timeout
	}
}

// WithDisable 禁用服务
func WithDisable() Option {
	return func(o *Server) {
		o.Status = StartStop
	}
}

// WithTrace 开启gin路由跟踪
func WithTrace() Option {
	return func(o *Server) {
		o.Trace = true
	}
}
