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

//WithScaleImage 压缩图片文件服务
func WithScaleImage(service string) Option {
	return func(a *NFS) {
		a.AllowScaleImage = true
		a.ScaleImageService = service
	}
}

//WithListFile 列出文件服务
func WithListFile(service string) Option {
	return func(a *NFS) {
		a.AllowListFile = true
		a.ListFileService = service
	}
}

//WithListDir 列出目录服务
func WithListDir(service string) Option {
	return func(a *NFS) {
		a.AllowListDir = true
		a.ListDirService = service
	}
}

//WithDownload 允许下载文件
func WithDownload(service string) Option {
	return func(a *NFS) {
		a.AllowDownload = true
		a.DownloadService = service
	}
}

//WithPreview 允许预览
func WithPreview(service string) Option {
	return func(a *NFS) {
		a.AllowPreview = true
		a.PreviewService = service
	}
}

//WithCreateDir 允许创建目录
func WithCreateDir(service string) Option {
	return func(a *NFS) {
		a.AllowCreateDir = true
		a.CreateDirService = service
	}
}

//WithRenameDir 允许创建目录
func WithRenameDir(service string) Option {
	return func(a *NFS) {
		a.AllowRenameDir = true
		a.RenameDirService = service
	}
}

//WithHWOBS 使用华为OBS
func WithHWOBS(ak string, sk string, bucket string, endpoint string) Option {
	return func(a *NFS) {
		a.CloudNFS = "OBS"
		a.AccessKey = ak
		a.SecretKey = sk
		a.Local = bucket
		a.Endpoint = endpoint
	}
}

//WithInclude 只包含目录
func WithIncludes(name ...string) Option {
	return func(a *NFS) {
		a.Includes = name
	}
}

//WithExcludes 排除目录
func WithExcludes(name ...string) Option {
	return func(a *NFS) {
		a.Excludes = name
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
