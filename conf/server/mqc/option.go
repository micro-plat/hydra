package mqc

//Option 配置选项
type Option func(*Server)

//WithTrace 构建api server配置信息
func WithTrace() Option {
	return func(a *Server) {
		a.Trace = true
	}
}

//WitchMasterSlave 设置为主备模式
func WitchMasterSlave() Option {
	return func(a *Server) {
		a.Sharding = 1
	}
}

//WitchSharding 设置为分片模式
func WitchSharding(i int) Option {
	return func(a *Server) {
		a.Sharding = i
	}
}

//WitchP2P 设置为对等模式
func WitchP2P() Option {
	return func(a *Server) {
		a.Sharding = 0
	}
}

//WithTimeout 构建api server配置信息
func WithTimeout(timeout int) Option {
	return func(a *Server) {
		a.Timeout = timeout
	}
}

//WithDisable 禁用任务
func WithDisable() Option {
	return func(a *Server) {
		a.Status = "stop"
	}
}

//WithEnable 启用任务
func WithEnable() Option {
	return func(a *Server) {
		a.Status = "start"
	}
}
