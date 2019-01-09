package http

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/engines"
	"github.com/micro-plat/lib4go/logger"
)

//WebResponsiveServer web 响应式服务器
type WebResponsiveServer struct {
	*ApiResponsiveServer
	webServer *WebServer
}

//NewWebResponsiveServer 构建基于注册中心的响应式web服务器
func NewWebResponsiveServer(registryAddr string, cnf conf.IServerConf, logger *logger.Logger) (h *WebResponsiveServer, err error) {
	h = &WebResponsiveServer{
		ApiResponsiveServer: &ApiResponsiveServer{registryAddr: registryAddr},
	}
	h.closeChan = make(chan struct{})
	h.currentConf = cnf
	h.Logger = logger
	h.pubs = make([]string, 0, 2)
	// 启动执行引擎
	h.engine, err = engines.NewServiceEngine(cnf, registryAddr, h.Logger)
	if err != nil {
		return nil, fmt.Errorf("%s:engine启动失败%v", cnf.GetServerName(), err)
	}
	if err = h.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return nil, err
	}
	if h.webServer, err = NewWebServer(cnf.GetServerName(),
		cnf.GetString("address", ":8080"),
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithLogger(logger),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithTimeout(cnf.GetInt("rTimeout", 10), cnf.GetInt("wTimeout", 10), cnf.GetInt("rhTimeout", 10))); err != nil {
		return
	}
	h.server = h.webServer
	if err = h.SetConf(true, h.currentConf); err != nil {
		return
	}
	return
}

//Restart 重启服务器
func (w *WebResponsiveServer) Restart(cnf conf.IServerConf) (err error) {
	w.Shutdown()
	time.Sleep(time.Second)
	w.closeChan = make(chan struct{})
	w.done = false
	w.once = sync.Once{}
	// 启动执行引擎
	w.engine, err = engines.NewServiceEngine(cnf, w.registryAddr, w.Logger)
	if err != nil {
		return fmt.Errorf("%s:engine启动失败%v", cnf.GetServerName(), err)
	}
	if err = w.engine.SetHandler(cnf.Get("__component_handler_").(component.IComponentHandler)); err != nil {
		return err
	}
	if w.server, err = NewWebServer(cnf.GetServerName(),
		cnf.GetString("address", ":8080"),
		nil,
		WithShowTrace(cnf.GetBool("trace", false)),
		WithLogger(w.Logger),
		WithName(cnf.GetPlatName(), cnf.GetSysName(), cnf.GetClusterName(), cnf.GetServerType()),
		WithTimeout(cnf.GetInt("rTimeout", 10), cnf.GetInt("wTimeout", 10), cnf.GetInt("rhTimeout", 10))); err != nil {
		return
	}

	if err = w.SetConf(true, cnf); err != nil {
		w.currentConf = cnf
		w.restarted = true
		return
	}
	return w.Start()
}
