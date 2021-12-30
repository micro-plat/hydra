package lnfs

import (
	"fmt"
	"net"
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
)

//currentModule 当前Module
var currentModule *Module

type LNFS struct {
	c         *nfs.NFS
	app       app.IAPPConf
	watcher   conf.IWatcher
	closch    chan struct{}
	isStarted bool
	*Module
	once     sync.Once
	services []string
}

func NewNFS(app app.IAPPConf, c *nfs.NFS) *LNFS {
	p, _ := app.GetProcessorConf()
	currentModule = newModule(c, p.ServicePrefix)
	return &LNFS{c: c, app: app, closch: make(chan struct{}), Module: currentModule, services: make([]string, 0, 3)}
}
func (c *LNFS) Start() error {
	if c.isStarted {
		return nil
	}
	c.isStarted = true
	go c.watch()
	return nil

}

//监控集群变化
func (c *LNFS) watch() {
	cluster, err := c.app.GetServerConf().GetCluster()
	if err != nil {
		return
	}
	c.watcher = cluster.Watch()
	notify := c.watcher.Notify()
	for {
		select {
		case <-notify:
			hosts, masterHost, currentAddr, isMaster, err := c.get()
			if err != nil {
				continue
			}
			c.Module.Update(hosts, masterHost, currentAddr, isMaster)

		case <-c.closch:
			c.watcher.Close()
			return
		}
	}
}

//get 从集群中获取数据
func (r *LNFS) get() (hosts []string, masterHost string, currentAddr string, isMaster bool, err error) {
	c, err := r.app.GetServerConf().GetCluster()
	if err != nil {
		return nil, "", "", false, err
	}
	if c.Current().GetPort() == "" {
		return nil, "", "", false, fmt.Errorf("系统未就绪")
	}
	hosts = make([]string, 0, 0)
	c.Iter(func(n conf.ICNode) bool {
		if !n.IsCurrent() {
			hosts = append(hosts, net.JoinHostPort(n.GetHost(), n.GetPort()))
		}
		if n.GetIndex() == 0 {
			masterHost = net.JoinHostPort(n.GetHost(), n.GetPort())
		}
		return true
	})
	isMaster = c.Current().GetIndex() == 0
	currentAddr = net.JoinHostPort(c.Current().GetHost(), c.Current().GetPort())
	return
}
func (r *LNFS) Close() error {
	r.once.Do(func() {
		close(r.closch)
		r.Module.Close()
	})
	return nil
}
