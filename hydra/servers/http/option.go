package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

type Handler interface {
	Handle(*gin.Context)
}
type option struct {
	readTimeout       int
	writeTimeout      int
	readHeaderTimeout int
	metric            *middleware.Metric
	serverType        string
	ginTrace          bool
}

//Option 配置选项
type Option func(*option)

//WithServerType 服务器类型
func WithServerType(t string) Option {
	return func(o *option) {
		o.serverType = t
	}
}

//WithTimeout 设置服务器超时时长
func WithTimeout(readTimeout int, writeTimeout int, readHeaderTimeout int) Option {
	return func(o *option) {
		o.readTimeout = readTimeout
		o.writeTimeout = writeTimeout
		o.readHeaderTimeout = readHeaderTimeout
	}
}

//WithGinTrace 是否启用gin注册跟踪
func WithGinTrace(b bool) Option {
	return func(o *option) {
		o.ginTrace = b
	}
}
