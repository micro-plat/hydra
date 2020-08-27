package apm

//Option 配置选项
type Option func(*APM)

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *APM) {
		a.m["enable"] = false
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *APM) {
		a.m["enable"] = true
	}
}

//WithConfig 配置名称
func WithConfig(config string) Option {
	return func(a *APM) {
		a.m["config"] = config
	}
}

//WithDB 配置DB参与监控
func WithDB(db bool) Option {
	return func(a *APM) {
		a.m["db"] = db
	}
}

//WithCache 配置Cache参与监控
func WithCache(cache bool) Option {
	return func(a *APM) {
		a.m["cache"] = cache
	}
}
