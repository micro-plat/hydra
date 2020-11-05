package gray

//Option 配置选项
type Option func(*Gray)

//WithFilter 设置筛选器
func WithFilter(filter string) Option {
	return func(a *Gray) {
		a.Filter = filter
	}
}

//WithUPCluster 设置上游集群名
func WithUPCluster(upcluster string) Option {
	return func(a *Gray) {
		a.UPCluster = upcluster
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *Gray) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *Gray) {
		a.Disable = false
	}
}
