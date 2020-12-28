package mqc

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	varqueue "github.com/micro-plat/hydra/conf/vars/queue"
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
	first    bool
}

//NewResponsive 创建响应式服务器
func NewResponsive(cnf app.IAPPConf) (h *Responsive, err error) {
	h = &Responsive{
		conf:     cnf,
		first:    true,
		log:      logger.New(cnf.GetServerConf().GetServerName()),
		pub:      pub.New(cnf.GetServerConf()),
		comparer: conf.NewComparer(cnf.GetServerConf(), mqc.MainConfName, mqc.SubConfName...),
	}
	app.Cache.Save(cnf)
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

	//发布集群节点
	if err = w.publish(); err != nil {
		err = fmt.Errorf("%s服务发布失败 %w", w.conf.GetServerConf().GetServerType(), err)
		w.Shutdown()
		return err
	}

	//监控服务节点变化并切换工作模式
	go w.watch()

	//动态监听任务
	services.MQC.Subscribe(func(t *queue.Queue) {
		if err := w.Server.Add(t); err != nil {
			w.log.Errorf("[%s]]服务[%v]添加失败 %w", w.conf.GetServerConf().GetServerType(), t, err)
		}
	})

	w.log.Infof("启动成功(%s,%s,[%d])", w.conf.GetServerConf().GetServerType(), w.Server.GetAddress(), w.Server.queues.Count())
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
			w.log.Info("mqc服务被禁用，不用重启")
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

//update 更新发布数据
func (w *Responsive) update(kv ...string) (err error) {
	addr := w.Server.GetAddress()
	serverName := strings.Split(addr, "://")[1]
	if err := w.pub.Update(serverName, addr, w.conf.GetServerConf().GetServerID(), kv...); err != nil {
		return err
	}

	return
}

//根据main.conf创建服务嚣
func (w *Responsive) getServer(cnf app.IAPPConf) (*Server, error) {
	nconf, err := cnf.GetMQCMainConf()
	if err != nil {
		return nil, fmt.Errorf("mqc服务器配置获取失败:%v", err)
	}
	proto, queuename, err := global.ParseProto(nconf.Addr)
	if err != nil {
		return nil, fmt.Errorf("mqc服务器地址有误:%w", err)
	}
	queueObj, err := cnf.GetMQCQueueConf()
	if err != nil {
		return nil, fmt.Errorf("mqc服务器监听队列配置有误:%w", err)
	}
	if global.IsLocal(proto) {
		return NewServer(proto, nil, queueObj.Queues...)
	}

	js, err := cnf.GetVarConf().GetConf(varqueue.TypeNodeName, queuename)
	if err != nil {
		return nil, fmt.Errorf("获取mqc服务器配置失败./var/%s/%s %w", varqueue.TypeNodeName, queuename, err)
	}
	return NewServer(proto, js.GetRaw(), queueObj.Queues...)
}

func init() {
	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return NewResponsive(c)
	}
	servers.Register(MQC, fn)
}

//MQC mqc服务器
const MQC = global.MQC
