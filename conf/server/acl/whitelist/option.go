package whitelist

import "github.com/micro-plat/hydra/conf"

//Option 配置选项
type Option func(*WhiteList)

// WithIPList WithIPList
func WithIPList(list ...*IPList) Option {
	return func(a *WhiteList) {
		for _, ip := range list {
			ip.ipm = conf.NewPathMatch(ip.IPS...)
			ip.rqm = conf.NewPathMatch(ip.Requests...)
			a.WhiteList = append(a.WhiteList, ip)
		}
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *WhiteList) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *WhiteList) {
		a.Disable = false
	}
}
