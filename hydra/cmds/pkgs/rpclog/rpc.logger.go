package rpclog

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

type layout struct {
	Level    string `json:"level" valid:"in(Off|Debug|Info|Warn|Error|Fatal|All),required"`
	Service  string `json:"service" valid:"required"`
	Interval string `json:"interval" valid:"required"`
}

type RPCLogger struct {
	registry registry.IRegistry
	pub      conf.IPub
	logger   *logger.Logger
	writer   *rpcWriter

	appenders []*RPCAppender

	layout *logger.Layout
}

//NewRPCLogger 创建RPC日志程序
func NewRPCLogger(registry registry.IRegistry, pub conf.IPub) (r *RPCLogger, err error) {
	r = &RPCLogger{
		pub:       pub,
		registry:  registry,
		appenders: make([]*RPCAppender, 0, 2),
		layout:    &logger.Layout{Type: "rpc", Level: "Info", Interval: "@every 3s"},
	}
	return r, nil

}

//MakeAppender 构建Appender
func (r *RPCLogger) MakeAppender(l *logger.Layout, event *logger.LogEvent) (logger.IAppender, error) {
	rpc, err := NewRPCAppender(r.writer, r.layout)
	if err != nil {
		return nil, err
	}
	r.appenders = append(r.appenders, rpc)
	return rpc, nil
}

//GetType 日志类型
func (r *RPCLogger) GetType() string {
	return "rpc"
}

//MakeUniq 获取日志标识
func (r *RPCLogger) MakeUniq(l *logger.Layout, event *logger.LogEvent) string {
	return "rpc"
}
func (r *RPCLogger) changed(c *conf.JSONConf) error {

	writer := newRPCWriter(setting.Service, r.platName, r.systemName, r.clusterName, r.serverTypes)
	r.writer = writer

	r.appender.Type = "rpc"
	r.appender.Level = setting.Level
	r.appender.Layout = c.GetString("layout")
	r.appender.Interval = setting.Interval

	logger.RegistryFactory(r, r.appender)

	return nil
}

//Close 关闭RPC日志
func (r *RPCLogger) Close() error {
	return nil
}
