package mqc

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/engines"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/logger"
)

//MqcResponsiveServer rpc 响应式服务器
type MqcResponsiveServer struct {
	server        *MqcServer
	engine        servers.IRegistryEngine
	registryAddr  string
	pubs          []string
	currentConf   conf.IServerConf
	closeChan     chan struct{}
	once          sync.Once
	done          bool
	shardingIndex int
	shardingCount int
	master        bool
	pubLock       sync.Mutex
	restarted     bool
	*logger.Logger
	mu sync.Mutex
}

//NewMqcResponsiveServer 创建mqc服务器
func NewMqcResponsiveServer(registryAddr string, cnf conf.IServerConf, logger *logger.Logger) (h *MqcResponsiveServer, err error) {
	h = &MqcResponsiveServer{
		closeChan:    make(chan struct{}),
		currentConf:  cnf,
		Logger:       logger,
		pubs:         make([]string, 0, 2),
		registryAddr: registryAddr,
	}
	// 启动执行引擎
	h.engine, err = engines.NewServiceEngine(cnf, registryAddr, h.Logger)
	if err != nil {
		return nil, fmt.Errorf("%s:engine启动失败%v", cnf.GetServerName(), err)
	}
	if err = h.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return nil, err
	}
	if h.server, err = NewMqcServer(cnf.GetServerName(),
		"",
		"",
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithLogger(logger)); err != nil {
		return
	}
	if err = h.SetConf(true, h.currentConf); err != nil {
		return
	}
	return
}

//Restart 重启服务器
func (w *MqcResponsiveServer) Restart(cnf conf.IServerConf) (err error) {
	w.Shutdown()
	time.Sleep(time.Second)
	w.done = false
	w.closeChan = make(chan struct{})
	w.once = sync.Once{}
	// 启动执行引擎
	w.engine, err = engines.NewServiceEngine(cnf, w.registryAddr, w.Logger)
	if err != nil {
		return fmt.Errorf("%s:engine启动失败%v", cnf.GetServerName(), err)
	}
	if err = w.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return err
	}
	if w.server, err = NewMqcServer(cnf.GetServerName(),
		"",
		"",
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithLogger(w.Logger)); err != nil {
		return
	}
	if err = w.SetConf(true, cnf); err != nil {
		return
	}
	if err = w.Start(); err == nil {
		w.currentConf = cnf
		w.restarted = true
		return
	}
	return err
}

//Start 启用服务
func (w *MqcResponsiveServer) Start() (err error) {
	if err = w.server.Run(); err != nil {
		return
	}
	return w.publish()
}

//Shutdown 关闭服务器
func (w *MqcResponsiveServer) Shutdown() {
	w.done = true
	w.once.Do(func() {
		close(w.closeChan)
	})
	w.unpublish()
	timeout := w.currentConf.GetInt("timeout", 10)
	w.server.Shutdown(time.Duration(timeout) * time.Second)
	if w.engine != nil {
		w.engine.Close()
	}
}

//GetAddress 获取服务器地址
func (w *MqcResponsiveServer) GetAddress() string {
	return w.server.GetAddress()
}

//GetStatus 获取当前服务器状态
func (w *MqcResponsiveServer) GetStatus() string {
	return w.server.GetStatus()
}

//GetServices 获取服务列表
func (w *MqcResponsiveServer) GetServices() []string {
	return w.engine.GetServices()
}

//Restarted 服务器是否已重启
func (w *MqcResponsiveServer) Restarted() bool {
	return w.restarted
}
