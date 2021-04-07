package nfs

//Option 配置选项
type Option func(*NFS)

//WithDisable 关闭
func WithDisable() Option {
	return func(a *NFS) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *NFS) {
		a.Disable = false
	}
}
