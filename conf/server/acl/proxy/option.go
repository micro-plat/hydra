package proxy

//Option 配置选项
type Option func(*Proxy)

//WithFilter 设置筛选器
func WithFilter(filter string) Option {
	return func(a *Proxy) {
		a.Filter = filter
	}
}

//WithUPCluster 设置上游集群名
func WithUPCluster(upcluster string) Option {
	return func(a *Proxy) {
		a.UPCluster = upcluster
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *Proxy) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *Proxy) {
		a.Disable = false
	}
}
