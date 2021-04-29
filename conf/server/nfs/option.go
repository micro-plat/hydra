package nfs

//Option 配置选项
type Option func(*NFS)

//WithRename 重命名文件名称
func WithRename() Option {
	return func(a *NFS) {
		a.Rename = true
	}
}

//WithWatch 启动文件夹监控
func WithWatch() Option {
	return func(a *NFS) {
		a.Watch = true
	}
}

//WithDisableUpload 禁用上传文件
func WithDisableUpload() Option {
	return func(a *NFS) {
		a.DiableUpload = true
	}
}

//WithUpload 设置上传服务名称
func WithUpload(service string) Option {
	return func(a *NFS) {
		a.UploadService = service
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
