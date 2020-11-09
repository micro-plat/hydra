package mqc

//Option 配置选项
type Option func(*Server)

//WithTrace 构建api server配置信息
func WithTrace() Option {
	return func(a *Server) {
		a.Trace = true
	}
}

//WithMasterSlave 设置为主备模式
func WithMasterSlave() Option {
	return func(a *Server) {
		a.Sharding = 1
	}
}

//WithSharding 设置为分片模式
func WithSharding(i int) Option {
	return func(a *Server) {
		a.Sharding = i
	}
}

//WithP2P 设置为对等模式
func WithP2P() Option {
	return func(a *Server) {
		a.Sharding = 0
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
