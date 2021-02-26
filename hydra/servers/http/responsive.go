package http

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"
)

//Responsive 响应式服务器
type Responsive struct {
	*Server
	conf     app.IAPPConf
	comparer conf.IComparer
	pub      pub.IPublisher
	log      logger.ILogger
}

//NewResponsive 创建响应式服务器
func NewResponsive(cnf app.IAPPConf) (h *Responsive, err error) {
	h = &Responsive{
		conf:     cnf,
		log:      logger.New(cnf.GetServerConf().GetServerName()),
		pub:      pub.New(cnf.GetServerConf()),
		comparer: conf.NewComparer(cnf.GetServerConf(), api.MainConfName, api.SubConfName...),
	}

	app.Cache.Save(cnf)
	if err := services.Def.DoSetup(cnf); err != nil {
		return nil, err
	}
	h.Server, err = h.getServer(cnf)
	return h, err
}

//Start 启用服务
func (w *Responsive) Start() (err error) {
	if err := services.Def.DoStarting(w.conf); err != nil {
		return err
	}

	if !w.conf.GetServerConf().IsStarted() {
		w.log.Warnf("%s被禁用，未启动", w.conf.GetServerConf().GetServerType())
		return
	}

	if err = w.Server.Start(); err != nil {
		err = fmt.Errorf("%s启动失败 %w", w.conf.GetServerConf().GetServerType(), err)
		return
	}

	if err = w.publish(); err != nil {
		err = fmt.Errorf("%s服务发布失败 %w", w.conf.GetServerConf().GetServerType(), err)
		w.Shutdown()
		return err
	}
	w.log.Infof("启动成功(%s,%s,[%d])", w.conf.GetServerConf().GetServerType(), w.Server.GetAddress(), w.serverNum())

	//服务启动成功后钩子
	if err := services.Def.DoStarted(w.conf); err != nil {
		err = fmt.Errorf("%s外部处理失败，关闭服务器 %w", w.conf.GetServerConf().GetServerType(), err)
		w.Shutdown()
		return err
	}
	return nil
}

//Notify 服务器配置变更通知
func (w *Responsive) Notify(c app.IAPPConf) (change bool, err error) {
	w.comparer.Update(c.GetServerConf())
	if !w.comparer.IsChanged() {
		return false, nil
	}
	if w.comparer.IsValueChanged() || w.comparer.IsSubConfChanged() {
		w.log.Info("关键配置发生变化，准备重启服务器")
		server, err := w.getServer(c)
		if err != nil {
			return false, err
		}

		w.Shutdown()
		w.conf = c
		app.Cache.Save(c)
		if !c.GetServerConf().IsStarted() {
			w.log.Info("api服务被禁用，不用重启")
			return true, nil
		}

		w.Server = server
		if err = w.Start(); err != nil {
			return false, err
		}
		return true, nil
	}
	app.Cache.Save(c)
	w.conf = c
	return true, nil
}

//Shutdown 关闭服务器
func (w *Responsive) Shutdown() {
	w.log.Infof("关闭[%s]服务...", w.conf.GetServerConf().GetServerType())
	w.Server.Shutdown()
	w.pub.Clear()
	if err := services.Def.DoClosing(w.conf); err != nil {
		w.log.Infof("关闭[%s]服务,出现错误", err)
		return
	}
	return
}

//publish 将当前服务器的节点信息发布到注册中心
func (w *Responsive) publish() (err error) {
	addr := w.Server.GetAddress()
	serverName := strings.Split(addr, "://")[1]

	if err := w.pub.Publish(serverName, addr, w.conf.GetServerConf().GetServerID()); err != nil {
		return err
	}

	return
}

//serverNum 获取服务数量
func (w *Responsive) serverNum() int {
	routers := w.Server.engine.Routes()
	serverMap := map[string]string{}
	for _, item := range routers {
		if _, ok := serverMap[item.Path]; !ok {
			serverMap[item.Path] = item.Path
		}
	}
	return len(serverMap)
}

//根据main.conf创建服务嚣
func (w *Responsive) getServer(cnf app.IAPPConf) (*Server, error) {
	tp := cnf.GetServerConf().GetServerType()
	apiConf, err := api.GetConf(cnf.GetServerConf())
	if err != nil {
		return nil, err
	}
	processorObj, err := cnf.GetProcessorConf()
	if err != nil {
		return nil, err
	}
	//从服务中获取路由
	sr := services.GetRouter(tp)
	routerconf, err := sr.GetRouters()
	if err != nil {
		return nil, err
	}
	routerlist := routerconf.GetRouters()
	processorObj.TreatRouters(routerlist)
	switch tp {
	case WS:
		return NewWSServer(tp,
			apiConf.GetWSAddress(),
			routerconf.GetRouters(),
			WithServerType(tp),
			WithTimeout(apiConf.GetRTimeout(), apiConf.GetWTimeout(), apiConf.GetRHTimeout()),
			WithGinTrace(apiConf.Trace),
		)
	case Web:
		return NewServer(tp,
			apiConf.GetWEBAddress(),
			routerconf.GetRouters(),
			WithServerType(tp),
			WithTimeout(apiConf.GetRTimeout(), apiConf.GetWTimeout(), apiConf.GetRHTimeout()),
			WithGinTrace(apiConf.Trace),
		)
	default:
		return NewServer(tp,
			apiConf.GetAPIAddress(),
			routerconf.GetRouters(),
			WithServerType(tp),
			WithTimeout(apiConf.GetRTimeout(), apiConf.GetWTimeout(), apiConf.GetRHTimeout()),
			WithGinTrace(apiConf.Trace),
		)
	}
}

func init() {
	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return NewResponsive(c)
	}
	servers.Register(API, fn)
	servers.Register(Web, fn)
	servers.Register(WS, fn)
}

//API api服务器
const API = global.API

//Web web服务器
const Web = global.Web

//WS web socket服务器
const WS = global.WS
