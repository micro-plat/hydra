package blacklist

//Option 配置选项
type Option func(*BlackList)

//WithIP 设置密钥
func WithIP(ip ...string) Option {
	return func(a *BlackList) {
		a.IPS = append(a.IPS, ip...)
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *BlackList) {
		a.EnableEncryption = true
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *BlackList) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *BlackList) {
		a.Disable = false
	}
}
