package rpclog

import (
	"fmt"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/rpc"
	"github.com/micro-plat/lib4go/logger"
)

type loggerSetting struct {
	Level    string `json:"level" valid:"in(Off|Debug|Info|Warn|Error|Fatal|All),required"`
	Service  string `json:"service" valid:"required"`
	Interval string `json:"interval" valid:"required"`
}

type RPCLogger struct {
	platName     string
	systemName   string
	serverTypes  []string
	clusterName  string
	rpcInvoker   *rpc.Invoker
	logger       *logger.Logger
	registryAddr string
	writer       *rpcWriter
	notify       chan *watcher.ContentChangeArgs
	watcher      *watcher.Watcher
	appenders    []*RPCAppender
	service      string
	appender     *logger.Appender
	currentConf  *conf.JSONConf
	closeChan    chan struct{}
	once         sync.Once
	lock         sync.RWMutex
}

//NewRPCLogger 创建RPC日志程序
func NewRPCLogger(spath string, registryAddr string, log *logger.Logger, platName string, systemName string, clusterName string, serverTypes []string) (r *RPCLogger, err error) {
	r = &RPCLogger{
		platName:     platName,
		systemName:   systemName,
		clusterName:  clusterName,
		serverTypes:  serverTypes,
		closeChan:    make(chan struct{}),
		logger:       log,
		registryAddr: registryAddr,
		appenders:    make([]*RPCAppender, 0, 2),
		appender:     &logger.Appender{Type: "rpc", Level: "Info", Interval: "@every 1m"},
	}
	//生成注册中心
	registry, err := registry.NewRegistryWithAddress(registryAddr, log)
	if err != nil {
		err = fmt.Errorf("初始化注册中心失败：%s:%v", registryAddr, err)
		return nil, err
	}

	//启动配置节点监控
	r.watcher = watcher.NewWatcher(spath, time.Second, registry, log)
	if r.notify, err = r.watcher.Start(); err != nil {
		return nil, err
	}
	go r.loopWatch()
	return r, nil

}

func (r *RPCLogger) loopWatch() {
	tkr := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-r.closeChan:
			return
		case <-tkr.C:
			tkr.Stop()
			if r.currentConf == nil {
				r.logger.Debug("未配置rpc日志")
			}
		case data := <-r.notify:
			c, err := conf.NewJSONConf(data.Content, data.Version)
			if err != nil {
				r.logger.Error(err)
				break
			}

			cmpr := conf.NewJSONComparer(r.currentConf, c)
			if cmpr.IsChanged() && cmpr.IsValueChanged("level", "layout", "interval", "service") {
				if err := r.changed(c); err != nil {
					r.logger.Error(err)
					break
				}
				r.currentConf = c
			}

		}
	}
}

//MakeAppender 构建Appender
func (r *RPCLogger) MakeAppender(l *logger.Appender, event *logger.LogEvent) (logger.IAppender, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	rpc, err := NewRPCAppender(r.writer, r.appender)
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
func (r *RPCLogger) MakeUniq(l *logger.Appender, event *logger.LogEvent) string {
	return "rpc"
}
func (r *RPCLogger) changed(c *conf.JSONConf) error {
	var setting loggerSetting
	if err := c.Unmarshal(&setting); err != nil {
		r.logger.Error(err)
		return err
	}
	if b, err := govalidator.ValidateStruct(&setting); !b {
		r.logger.Error(fmt.Errorf("rpc logger配置有误:%v", err))
		return err
	}

	if _, err := time.ParseDuration(setting.Interval); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.service != setting.Service {
		_, domain, server, err := rpc.ResolvePath(setting.Service, "", "")
		if err != nil || domain == "" || server == "" {
			return fmt.Errorf("%s不合法 %v", setting.Service, err)
		}
		if r.rpcInvoker != nil {
			r.rpcInvoker.Close()
		}
		r.rpcInvoker = rpc.NewInvoker(domain, server, r.registryAddr)
		r.service = setting.Service
	}

	writer := newRPCWriter(setting.Service, r.rpcInvoker, r.platName, r.systemName, r.clusterName, r.serverTypes)
	r.writer = writer

	r.appender.Type = "rpc"
	r.appender.Level = setting.Level
	r.appender.Layout = c.GetString("layout")
	r.appender.Interval = setting.Interval

	for _, app := range r.appenders {
		app.Reset(setting.Interval, writer)
	}
	r.logger.Debug("rpc 日志配置成功")
	r.once.Do(func() {
		logger.RegistryFactory(r, r.appender)
	})
	return nil
}

//Close 关闭RPC日志
func (r *RPCLogger) Close() error {
	close(r.closeChan)
	if r.rpcInvoker != nil {
		r.rpcInvoker.Close()
	}
	if r.watcher != nil {
		r.watcher.Close()
	}
	return nil
}
