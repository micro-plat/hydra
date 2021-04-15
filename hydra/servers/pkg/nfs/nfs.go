package nfs

import (
	"fmt"
	"net"
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/concurrent/cmap"
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
	return &cnfs{c: c, app: app, closch: make(chan struct{}), module: newModule(c)}
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
			c.module.Update(hosts, masterHost, currentAddr, isMaster)

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
	isMaster = c.Current().GetIndex() == 0
	currentAddr = net.JoinHostPort(c.Current().GetHost(), c.Current().GetPort())
	return
}
func (r *cnfs) Close() error {
	r.once.Do(func() {
		close(r.closch)
		r.module.Close()
	})
	return nil
}

func init() {
	global.OnReady(func() {
		//处理服务初始化
		services.Def.OnSetup(func(c app.IAPPConf) error {
			//取消服务注册
			unRegistry(c.GetServerConf().GetServerType())
			closeNFS(c.GetServerConf().GetServerType())
			n, err := c.GetNFSConf()
			if err != nil {
				return err
			}
			if n.Disable {
				return nil
			}

			//构建并缓存nfs
			cnfs := newNFS(c, n)
			nfsCaches.Set(c.GetServerConf().GetServerType(), cnfs)

			//注册服务
			registry(c.GetServerConf().GetServerType(), cnfs, n)
			return nil
		})

		//处理服务启动完成
		services.Def.OnStarted(func(c app.IAPPConf) error {
			return startNFS(c.GetServerConf().GetServerType())
		})

	})

}

var nfsCaches cmap.ConcurrentMap = cmap.New(2)

func startNFS(tp string) error {
	v, ok := nfsCaches.Get(tp)
	if !ok {
		return nil
	}
	m := v.(*cnfs)
	return m.Start()
}

func closeNFS(tp string) error {
	nfsCaches.RemoveIterCb(func(k string, v interface{}) bool {
		if k == tp {
			m := v.(*cnfs)
			m.Close()
			return true
		}
		return false
	})
	return nil

}
func unRegistry(tp string) {
	if tp == global.API || tp == global.Web {
		//注册服务
		services.Def.Remove(SVSUpload, tp)
		services.Def.Remove(SVSDonwload, tp)

		//内部服务
		services.Def.Remove(rmt_fp_get, tp)
		services.Def.Remove(rmt_fp_notify, tp)
		services.Def.Remove(rmt_fp_query, tp)
		services.Def.Remove(rmt_file_download, tp)
	}
}
func registry(tp string, cnfs *cnfs, cnf *nfs.NFS) {
	if tp == global.API {
		//注册服务
		if !cnf.DiableUpload {
			services.Def.API(SVSUpload, cnfs.Upload)
		}

		if cnf.AllowDownload {
			services.Def.API(SVSDonwload, cnfs.Download)
		}

		//内部服务
		services.Def.API(rmt_fp_get, cnfs.GetFP)
		services.Def.API(rmt_fp_notify, cnfs.RecvNotify)
		services.Def.API(rmt_fp_query, cnfs.Query)
		services.Def.API(rmt_file_download, cnfs.GetFile)
	}

	if tp == global.Web {

		if !cnf.DiableUpload {
			services.Def.Web(SVSUpload, cnfs.Upload)
		}

		if cnf.AllowDownload {
			services.Def.Web(SVSDonwload, cnfs.Download)
		}

		//内部服务
		services.Def.Web(rmt_fp_get, cnfs.GetFP)
		services.Def.Web(rmt_fp_notify, cnfs.RecvNotify)
		services.Def.Web(rmt_fp_query, cnfs.Query)
		services.Def.Web(rmt_file_download, cnfs.GetFile)
	}
}
