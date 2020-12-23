package cli


type option struct {
	version string
	usage string
}

//Option 配置选项
type Option func(*option)

//WithVersion 设置版本号
func WithVersion(version string) Option {
	return func(o *option) {
		o.version=version
	}
}


//WithUsage 设置使用说明
func WithUsage(usage string) Option {
	return func(o *option) {
		o.usage=usage
	}
}