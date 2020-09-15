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

//WithName 配置名称
func WithName(name string) Option {
	return func(a *APM) {
		a.Name = name
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
