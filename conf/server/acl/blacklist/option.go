package blacklist

//Option 配置选项
type Option func(*BlackList)

//WithIP 设置密钥
func WithIP(ip ...string) Option {
	return func(a *BlackList) {
		a.IPS = append(a.IPS, ip...)
	}
}
