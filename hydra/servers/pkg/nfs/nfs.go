package nfs

import (
	"fmt"
	"net"
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
	"github.com/micro-plat/hydra/global"
)

type cnfs struct {
	c         *nfs.NFS
	app       app.IAPPConf
	watcher   conf.IWatcher
	closch    chan struct{}
	isStarted bool
	module    *module
	once      sync.Once
}

func newNFS(app app.IAPPConf, c *nfs.NFS) *cnfs {
	return &cnfs{c: c, app: app, closch: make(chan struct{}), module: newModule(c.Local)}
}
func (c *cnfs) Start() error {
	if c.isStarted {
		return nil
	}
	c.isStarted = true
	go c.watch()
	return nil

}

//监控集群变化
func (c *cnfs) watch() {
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
			trace(fmt.Sprintf("change:hosts:%v,master:%s,current:%s,isMaster:%v", hosts, masterHost, currentAddr, isMaster))
			c.module.Update(c.c.Local, hosts, currentAddr, masterHost, isMaster)

		case <-c.closch:
			c.watcher.Close()
			return
		}
	}
}

//get 从集群中获取数据
func (r *cnfs) get() (hosts []string, masterHost string, currentAddr string, isMaster bool, err error) {
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

	//是否是主集群
	isMaster = c.Current().GetIndex() == 0
	currentAddr = net.JoinHostPort(c.Current().GetHost(), c.Current().GetPort())
	return
}
func (r *cnfs) Close() error {
	r.once.Do(func() {
		close(r.closch)
		if r.module != nil {
			r.module.Close()
		}

	})
	return nil
}

func trace(msg ...interface{}) {
	global.Def.Log().Debug(msg...)
}
