package rpc

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

//RpcResponsiveServer rpc 响应式服务器
type RpcResponsiveServer struct {
	server       *RpcServer
	engine       servers.IRegistryEngine
	registryAddr string
	pubs         []string
	currentConf  conf.IServerConf
	closeChan    chan struct{}
	once         sync.Once
	done         bool
	pubLock      sync.Mutex
	restarted    bool
	*logger.Logger
	mu sync.Mutex
}

//NewRpcResponsiveServer 创建rpc服务器
func NewRpcResponsiveServer(registryAddr string, cnf conf.IServerConf, logger *logger.Logger) (h *RpcResponsiveServer, err error) {
	h = &RpcResponsiveServer{
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
	if h.server, err = NewRpcServer(cnf.GetServerName(),
		cnf.GetString("address", "8081"),
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
func (w *RpcResponsiveServer) Restart(cnf conf.IServerConf) (err error) {
	w.Shutdown()
	time.Sleep(time.Second)
	w.done = false
	w.closeChan = make(chan struct{})
	w.currentConf = cnf
	w.once = sync.Once{}
	// 启动执行引擎
	w.engine, err = engines.NewServiceEngine(cnf, w.registryAddr, w.Logger)
	if err != nil {
		return fmt.Errorf("%s:engine启动失败%v", cnf.GetServerName(), err)
	}
	if err = w.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return err
	}
	if w.server, err = NewRpcServer(cnf.GetServerName(),
		cnf.GetString("address", "8080"),
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
func (w *RpcResponsiveServer) Start() (err error) {
	if err = w.server.Run(); err != nil {
		return
	}
	return w.publish()
}

//Shutdown 关闭服务器
func (w *RpcResponsiveServer) Shutdown() {
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
func (w *RpcResponsiveServer) GetAddress() string {
	return w.server.GetAddress()
}

//GetStatus 获取当前服务器状态
func (w *RpcResponsiveServer) GetStatus() string {
	return w.server.GetStatus()
}

//GetServices 获取服务列表
func (w *RpcResponsiveServer) GetServices() []string {
	svs := w.engine.GetServices()
	nsevice := make([]string, 0, len(svs))
	for _, sv := range svs {
		if w.server.Find(sv) {
			nsevice = append(nsevice, sv)
		}
	}
	//servers.Trace(w.Infof, w.currentConf.GetServerName(), "发布服务：", nsevice)
	return nsevice
}

//Restarted 服务器是否已重启
func (w *RpcResponsiveServer) Restarted() bool {
	return w.restarted
}
