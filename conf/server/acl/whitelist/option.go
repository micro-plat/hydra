package whitelist

//Option 配置选项
type Option func(*IPList)

//NewIPList 构建IPLIST
func NewIPList(request string, opts ...Option) *IPList {
	ip := &IPList{
		Requests: []string{request},
	}
	for _, opt := range opts {
		opt(ip)
	}
	return ip
}

//WithIP 设置密钥
func WithIP(ip ...string) Option {
	return func(a *IPList) {
		a.IPS = append(a.IPS, ip...)
	}
}
