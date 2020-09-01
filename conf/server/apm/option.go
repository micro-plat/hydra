package apm

//Option 配置选项
type Option func(*APM)

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *APM) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *APM) {
		a.Disable = false
	}
}

//WithConfig 配置名称
func WithConfig(config string) Option {
	return func(a *APM) {
		a.Config = config
	}
}

//WithDB 配置DB参与监控
func WithDB(db ...string) Option {
	return func(a *APM) {
		a.DB = db
	}
}

//WithCache 配置Cache参与监控
func WithCache(cache ...string) Option {
	return func(a *APM) {
		a.Cache = cache
	}
}
