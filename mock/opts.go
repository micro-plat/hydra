package mock

import (
	"net/http"

	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
)

//Option 配置选项
type Option func(*mock)

//WithServerType 服务器类型
func WithServerType(t string) Option {
	return func(o *mock) {
		global.Def.ServerTypes = []string{t}
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

//WithEncoding 设置编码格式
func WithEncoding(encoding string) Option {
	return func(o *mock) {
		o.encoding = encoding
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

//WithPlatName 设置平台名称
func WithPlatName(platName string) Option {
	return func(o *mock) {
		global.Def.PlatName = platName
	}
}

//WithSystemName 设置系统名称
func WithSystemName(sysName string) Option {
	return func(o *mock) {
		global.Def.SysName = sysName
	}
}

//WithClusterName 设置集群名称
func WithClusterName(clusterName string) Option {
	return func(o *mock) {
		global.Def.ClusterName = clusterName
	}
}

//WithRegistry 设置注册中心地址
func WithRegistry(addr string) Option {
	return func(o *mock) {
		global.Def.RegistryAddr = addr
	}
}

//WithConf 设置配置对象
func WithConf(conf creator.IConf) Option {
	return func(o *mock) {
		o.Conf = conf
	}
}

//WithRequest 设置HTTP请求
func WithRequest(request *http.Request) Option {
	return func(o *mock) {
		o.Request = request
	}
}

//WithResponse 设置HTTP writer
func WithResponse(response http.ResponseWriter) Option {
	return func(o *mock) {
		o.Response = response
	}
}
