package wss

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	wssconf "github.com/micro-plat/hydra/conf/server/wss"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"
)

type Responsive struct {
	conf     app.IAPPConf
	comparer conf.IComparer
	pub      pub.IPublisher
	log      logger.ILogger
	server   *Server
	client   *Client
}

func NewResponsive(cnf app.IAPPConf) (*Responsive, error) {
	r := &Responsive{
		conf:     cnf,
		log:      logger.New(cnf.GetServerConf().GetServerName()),
		pub:      pub.New(cnf.GetServerConf()),
		comparer: conf.NewComparer(cnf.GetServerConf(), wssconf.MainConfName, wssconf.SubConfName...),
	}
	app.Cache.Save(cnf)
	if err := services.Def.DoSetup(cnf); err != nil {
		return nil, err
	}
	if err := r.build(cnf); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Responsive) build(cnf app.IAPPConf) error {
	main, err := wssconf.GetConf(cnf.GetServerConf())
	if err != nil {
		return err
	}
	processorObj, err := cnf.GetProcessorConf()
	if err != nil {
		return err
	}
	routerObj, err := services.GetRouter(cnf.GetServerConf().GetServerType()).BuildRouters(processorObj.ServicePrefix)
	if err != nil {
		return err
	}
	engine := newDispatcher(cnf.GetServerConf().GetServerType(), routerObj.GetRouters())
	switch cnf.GetServerConf().GetServerType() {
	case global.WSSServer:
		routes, err := wssconf.GetRoutes(cnf.GetServerConf())
		if err != nil {
			return err
		}
		r.server = NewServer(main, engine, routes.Routes, r.log)
		r.client = nil
	case global.WSSClient:
		r.client = NewClient(main, engine, r.log)
		r.server = nil
	default:
		return fmt.Errorf("unsupported wss server type %s", cnf.GetServerConf().GetServerType())
	}
	return nil
}

func (r *Responsive) Start() error {
	if err := services.Def.DoStarting(r.conf); err != nil {
		return err
	}
	if !r.conf.GetServerConf().IsStarted() {
		r.log.Warnf("%s被禁用，未启动", r.conf.GetServerConf().GetServerType())
		return nil
	}
	if r.server != nil {
		if err := r.server.Start(); err != nil {
			return err
		}
		if err := r.publish(); err != nil {
			r.Shutdown()
			return err
		}
		r.log.Infof("启动成功(%s,%s,[%d])", r.conf.GetServerConf().GetServerType(), r.server.GetAddress(), r.server.ServiceNum())
	}
	if r.client != nil {
		if err := r.client.Start(); err != nil {
			return err
		}
		r.log.Infof("启动成功(%s,%s,[%d])", r.conf.GetServerConf().GetServerType(), r.client.GetAddress(), r.client.ServiceNum())
	}
	return services.Def.DoStarted(r.conf)
}

func (r *Responsive) Notify(cnf app.IAPPConf) (bool, error) {
	r.comparer.Update(cnf.GetServerConf())
	if !r.comparer.IsChanged() {
		return false, nil
	}
	if err := services.Def.DoSetup(cnf); err != nil {
		return false, err
	}
	r.Shutdown()
	r.conf = cnf
	app.Cache.Save(cnf)
	if err := r.build(cnf); err != nil {
		return false, err
	}
	if cnf.GetServerConf().IsStarted() {
		if err := r.Start(); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *Responsive) Shutdown() {
	r.log.Infof("关闭[%s]服务...", r.conf.GetServerConf().GetServerType())
	if r.server != nil {
		r.server.Shutdown()
		r.pub.Clear()
	}
	if r.client != nil {
		r.client.Shutdown()
	}
	if err := services.Def.DoClosing(r.conf); err != nil {
		r.log.Infof("关闭[%s]服务,出现错误", err)
	}
}

func (r *Responsive) publish() error {
	if r.server == nil {
		return nil
	}
	addr := r.server.GetAddress()
	serverName := strings.TrimPrefix(addr, "ws://")
	return r.pub.Publish(serverName, addr, r.conf.GetServerConf().GetServerID())
}

func init() {
	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return NewResponsive(c)
	}
	servers.Register(global.WSSServer, fn)
	servers.Register(global.WSSClient, fn)
}
