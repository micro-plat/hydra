package api

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
)

type Handler interface {
	Handle(*gin.Context)
}
type option struct {
	*logger.Logger
	readTimeout       int
	writeTimeout      int
	readHeaderTimeout int
	showTrace         bool
	metric            *middleware.Metric
	serverType        string
	tls               []string
}

//Option 配置选项
type Option func(*option)

//WithServerType 服务器类型
func WithServerType(t string) Option {
	return func(o *option) {
		o.serverType = t
	}
}

//WithLogger 设置日志记录组件
func WithLogger(logger *logger.Logger) Option {
	return func(o *option) {
		o.Logger = logger
	}
}

//WithShowTrace 显示跟踪信息
func WithShowTrace(b bool) Option {
	return func(o *option) {
		o.showTrace = b
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

//WithTLS 设置TLS证书(pem,key)
func WithTLS(tls []string) Option {
	return func(o *option) {
		if len(tls) == 2 {
			o.tls = tls
		}
	}
}
