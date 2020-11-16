package rpc

//Option 配置选项
type Option func(*Server)

//WithTrace 构建api server配置信息
func WithTrace() Option {
	return func(a *Server) {
		a.Trace = true
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
func WithDNS(host string, ip ...string) Option {
	return func(a *Server) {
		a.Host = host
		if len(ip) > 0 {
			a.Domain = ip[0]
		}
	}
}

//WithMaxRecvMsgSize 最大接收字节数
func WithMaxRecvMsgSize(maxRecvMsgSize int) Option {
	return func(a *Server) {
		a.MaxRecvMsgSize = maxRecvMsgSize
	}
}

//WithMaxSendMsgSize 最大发送字节数
func WithMaxSendMsgSize(maxRecvMsgSize int) Option {
	return func(a *Server) {
		a.MaxRecvMsgSize = maxRecvMsgSize
	}
}
