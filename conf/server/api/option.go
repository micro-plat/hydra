package api

//Option 配置选项
type Option func(*Server)

//WithTrace 构建api server配置信息
func WithTrace() Option {
	return func(a *Server) {
		a.Trace = true
	}
}

//WithTimeout 构建api server配置信息
func WithTimeout(rtimeout int, wtimout int) Option {
	return func(a *Server) {
		a.RTimeout = rtimeout
		a.WTimeout = wtimout
	}
}

//WithHeaderReadTimeout 构建api server配置信息
func WithHeaderReadTimeout(htimeout int) Option {
	return func(a *Server) {
		a.RHTimeout = htimeout
	}
}

//WithDisable 禁用任务
func WithDisable() Option {
	return func(a *Server) {
		a.Status = StartStop
	}
}

//WithEnable 启用任务
func WithEnable() Option {
	return func(a *Server) {
		a.Status = StartStatus
	}
}

//WithDNS 设置请求域名
func WithDNS(ip ...string) Option {
	return func(a *Server) {
		if len(ip) > 0 {
			a.Domain = ip[0]
		}
	}
}
