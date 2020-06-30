package cron

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"
)

//Responsive 响应式服务器
type Responsive struct {
	*Server
	conf     server.IServerConf
	comparer conf.IComparer
	pub      pub.IPublisher
	log      logger.ILogger
	first    bool
}

//NewResponsive 创建响应式服务器
func NewResponsive(cnf server.IServerConf) (h *Responsive, err error) {
	h = &Responsive{
		conf:     cnf,
		first:    true,
		log:      logger.New(cnf.GetMainConf().GetServerName()),
		pub:      pub.New(cnf.GetMainConf()),
		comparer: conf.NewComparer(cnf.GetMainConf(), cron.MainConfName, cron.SubConfName...),
	}
	server.Cache.Save(cnf)
	h.Server, err = h.getServer(cnf)
	return h, err
}

//Start 启用服务
func (w *Responsive) Start() (err error) {
	if err := services.Def.DoStarting(w.conf); err != nil {
		return err
	}
	if err = w.Server.Start(); err != nil {
		err = fmt.Errorf("启动失败 %w", err)
		return
	}

	//发布集群节点
	if err = w.publish(); err != nil {
		err = fmt.Errorf("服务发布失败 %w", err)
		w.Shutdown()
		return err
	}

	//监控服务节点变化并切换工作模式
	go w.watch()

	w.subscribe()

	w.log.Infof("启动成功(%s,%s,%d)", w.conf.GetMainConf().GetServerType(), w.Server.GetAddress(), len(w.conf.GetCRONTaskConf().Tasks))
	return nil
}

//Notify 服务器配置变更通知
func (w *Responsive) Notify(c server.IServerConf) (change bool, err error) {
	w.comparer.Update(c.GetMainConf())
	if !w.comparer.IsChanged() {
		return false, nil
	}
	if w.comparer.IsValueChanged() || w.comparer.IsSubConfChanged() {
		w.log.Info("关键配置发生变化，准备重启服务器")
		w.Shutdown()

		server.Cache.Save(c)
		w.Server, err = w.getServer(c)
		if err != nil {
			return false, err
		}
		if err = w.Start(); err != nil {
			return false, err
		}
		w.conf = c
		return true, nil
	}
	server.Cache.Save(c)
	w.conf = c
	return true, nil
}

func (w *Responsive) subscribe() {
	//动态监听任务
	services.CRON.Subscribe(func(t *task.Task) {
		if err := w.Server.Add(t); err != nil {
			w.log.Errorf("服务[%v]添加失败 %w", t, err)
		}
	})
}

//Shutdown 关闭服务器
func (w *Responsive) Shutdown() {
	w.log.Infof("关闭[%s]服务...", w.conf.GetMainConf().GetServerType())
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
	if err := w.pub.Publish(serverName, addr, w.conf.GetMainConf().GetServerID()); err != nil {
		return err
	}
	return
}

//update 更新发布数据
func (w *Responsive) update(kv ...string) (err error) {
	addr := w.Server.GetAddress()
	serverName := strings.Split(addr, "://")[1]
	if err := w.pub.Update(serverName, addr, w.conf.GetMainConf().GetServerID(), kv...); err != nil {
		return err
	}

	return
}

//根据main.conf创建服务嚣
func (w *Responsive) getServer(cnf server.IServerConf) (*Server, error) {
	//初始化server
	return NewServer(cnf.GetCRONTaskConf().Tasks...)
}

func init() {
	fn := func(c server.IServerConf) (servers.IResponsiveServer, error) {
		return NewResponsive(c)
	}
	servers.Register(CRON, fn)
}

//CRON cron服务器
const CRON = global.CRON
