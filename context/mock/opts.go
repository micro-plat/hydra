package mock

import "github.com/micro-plat/lib4go/types"

//Option 配置选项
type Option func(*mock)

//WithServerType 服务器类型
func WithServerType(t string) Option {
	return func(o *mock) {
		o.serverType = t
	}
}

//WithURL 设置URL参数
func WithURL(url string) Option {
	return func(o *mock) {
		o.URL = url
	}
}

//WithService 设置service参数
func WithService(service string) Option {
	return func(o *mock) {
		o.Service = service
	}
}

//WithRHeaders 设置header参数
func WithRHeaders(header types.XMap) Option {
	return func(o *mock) {
		o.RHeaders = header
	}
}

//WithCookies 设置cookie参数
func WithCookies(cookie types.XMap) Option {
	return func(o *mock) {
		o.Cookies = cookie
	}
}
