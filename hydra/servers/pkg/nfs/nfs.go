package nfs

import (
	"net"
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
)

type cnfs struct {
	c       *nfs.NFS
	app     app.IAPPConf
	watcher conf.IWatcher
	closch  chan struct{}
	module  *module
	once    sync.Once
}

func newNFS(c *nfs.NFS) *cnfs {
	return &cnfs{c: c}
}
func (c *cnfs) Start() error {
	hosts, masterHost, isMaster, err := c.get()
	if err != nil {
		return err
	}
	c.module, err = newModule(c.c.Local, hosts, masterHost, isMaster)
	if err != nil {
		return err
	}
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
			hosts, masterHost, isMaster, err := c.get()
			if err != nil {
				return
			}
			c.module.Update(hosts, c.c.Local, masterHost, isMaster)
		case <-c.closch:
			c.watcher.Close()
			return
		}
	}
}

//get 从集群中获取数据
func (r *cnfs) get() (hosts []string, masterHost string, isMaster bool, err error) {
	c, err := r.app.GetServerConf().GetCluster()
	if err != nil {
		return nil, "", false, err
	}
	hosts = make([]string, 0, 0)
	c.Iter(func(n conf.ICNode) bool {
		if !n.IsCurrent() {
			hosts = append(hosts, net.JoinHostPort(n.GetHost(), n.GetPort()))
		}
		if n.IsMaster(n.GetIndex()) {
			masterHost = net.JoinHostPort(n.GetHost(), n.GetPort())
		}
		return true
	})

	//是否是主集群
	isMaster = c.Current().IsMaster(c.Current().GetIndex())
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
