package cron

import (
	"github.com/micro-plat/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/logger"
)

type option struct {
	ip string
	*logger.Logger
	metric    *middleware.Metric
	showTrace bool
}

//Option 配置选项
type Option func(*option)

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

//WithIP 设置ip地址
func WithIP(ip string) Option {
	return func(o *option) {
		o.ip = ip
	}
}

//WithMetric 设置基于influxdb的系统监控组件
func WithMetric(host string, dataBase string, userName string, password string, cron string) Option {
	return func(o *option) {
		o.metric.Restart(host, dataBase, userName, password, cron, o.Logger)
	}
}
