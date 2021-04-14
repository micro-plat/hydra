package nfs

//Option 配置选项
type Option func(*NFS)

//WithDisableUpload 允许下载文件
func WithDisableUpload() Option {
	return func(a *NFS) {
		a.DiableUpload = true
	}
}

//WithAllowDownload 允许下载文件
func WithAllowDownload() Option {
	return func(a *NFS) {
		a.AllowDownload = true
	}
}

//WithDomain 下载域名
func WithDomain(domain string) Option {
	return func(a *NFS) {
		a.Domain = domain
	}
}

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
