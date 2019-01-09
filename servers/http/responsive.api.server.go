package http

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

type IServer interface {
	Run() error
	Shutdown(timeout time.Duration)
	GetStatus() string
	GetAddress() string
	CloseCircuitBreaker() error
	SetCircuitBreaker(*conf.CircuitBreaker) error

	SetRouters(routers []*conf.Router) (err error)
	SetJWT(auth *conf.Auth) error
	SetAjaxRequest(allow bool) error
	SetHosts(conf.Hosts) error
	SetStatic(*conf.Static) error
	SetMetric(*conf.Metric) error
	SetHeader(conf.Headers) error
	StopMetric() error
}

//ApiResponsiveServer api 响应式服务器
type ApiResponsiveServer struct {
	server       IServer
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

//NewApiResponsiveServer 创建API服务器
func NewApiResponsiveServer(registryAddr string, cnf conf.IServerConf, logger *logger.Logger) (h *ApiResponsiveServer, err error) {
	h = &ApiResponsiveServer{
		closeChan:    make(chan struct{}),
		currentConf:  cnf,
		Logger:       logger,
		pubs:         make([]string, 0, 2),
		registryAddr: registryAddr,
	}
	// 启动执行引擎
	h.engine, err = engines.NewServiceEngine(cnf, registryAddr, h.Logger)
	if err != nil {
		return nil, fmt.Errorf("engine启动失败%v", err)
	}
	if err = h.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return nil, err
	}
	if h.server, err = NewApiServer(cnf.GetServerName(),
		cnf.GetString("address", ":8090"),
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithLogger(logger),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithTimeout(cnf.GetInt("rTimeout", 10), cnf.GetInt("wTimeout", 10), cnf.GetInt("rhTimeout", 10))); err != nil {
		return
	}
	if err = h.SetConf(true, h.currentConf); err != nil {
		return
	}

	return
}

//Restart 重启服务器
func (w *ApiResponsiveServer) Restart(cnf conf.IServerConf) (err error) {
	w.Shutdown()
	time.Sleep(time.Second)
	w.done = false
	w.closeChan = make(chan struct{})
	w.once = sync.Once{}
	// 启动执行引擎
	w.engine, err = engines.NewServiceEngine(cnf, w.registryAddr, w.Logger)
	if err != nil {
		return fmt.Errorf("engine启动失败%v", err)
	}
	if err = w.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return err
	}
	if w.server, err = NewApiServer(cnf.GetServerName(),
		cnf.GetString("address", ":8090"),
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithLogger(w.Logger),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithTimeout(cnf.GetInt("rTimeout", 10), cnf.GetInt("wTimeout", 10), cnf.GetInt("rhTimeout", 10))); err != nil {
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
func (w *ApiResponsiveServer) Start() (err error) {
	if err = w.server.Run(); err != nil {
		return
	}
	if err = w.publish(); err != nil {
		w.Shutdown()
		return err
	}
	return nil
}

//Shutdown 关闭服务器
func (w *ApiResponsiveServer) Shutdown() {
	w.done = true
	w.once.Do(func() {
		close(w.closeChan)
	})
	w.unpublish()
	w.server.Shutdown(10 * time.Second)
	if w.engine != nil {
		w.engine.Close()
	}
}

//GetAddress 获取服务器地址
func (w *ApiResponsiveServer) GetAddress() string {
	return w.server.GetAddress()
}

//GetStatus 获取当前服务器状态
func (w *ApiResponsiveServer) GetStatus() string {
	return w.server.GetStatus()
}

//GetServices 获取服务列表
func (w *ApiResponsiveServer) GetServices() []string {
	return w.engine.GetServices()
}

//Restarted 服务器是否已重启
func (w *ApiResponsiveServer) Restarted() bool {
	return w.restarted
}
