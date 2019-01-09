package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/servers/http/middleware"
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
	platName          string
	systemName        string
	clusterName       string
	serverType        string
}

//Option 配置选项
type Option func(*option)

func WithName(platName string, systemName string, clusterName string, serverType string) Option {
	return func(o *option) {
		o.platName = platName
		o.systemName = systemName
		o.clusterName = clusterName
		o.serverType = serverType
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

//WithMetric 设置基于influxdb的系统监控组件
func WithMetric(host string, dataBase string, userName string, password string, cron string) Option {
	return func(o *option) {
		o.metric.Restart(host, dataBase, userName, password, cron, o.Logger)
	}
}
