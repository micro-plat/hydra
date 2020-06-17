package server

import (
	"sort"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
)

var cluster = cmap.New(2)

func getCluster(pub conf.IPub, rgst registry.IRegistry) (s *Cluster, err error) {
	_, c, err := cluster.SetIfAbsentCb(pub.GetServerPubPath(), func(...interface{}) (interface{}, error) {
		return NewCluster(pub, rgst)
	})
	if err != nil {
		return nil, err
	}
	return c.(*Cluster), nil
}

//Cluster 集群
type Cluster struct {
	conf.IPub
	registry registry.IRegistry
	current  conf.ICNode
	nodes    cmap.ConcurrentMap
	closeCh  chan struct{}
	lock     sync.Mutex
	watchers cmap.ConcurrentMap
}

//NewCluster 管理服务器的主配置信息
func NewCluster(pub conf.IPub, rgst registry.IRegistry) (s *Cluster, err error) {
	s = &Cluster{
		IPub:     pub,
		registry: rgst,
		nodes:    cmap.New(4),
		watchers: cmap.New(2),
		closeCh:  make(chan struct{}),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//Current 获取当前节点
func (c *Cluster) Current() conf.ICNode {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.current == nil {
		return &CNode{}
	}
	return c.current.Clone()
}

//Iter 迭代所有集群节点
func (c *Cluster) Iter(f func(conf.ICNode) bool) {
	nodes := c.nodes.Items()
	for _, c := range nodes {
		node := (c.(conf.ICNode)).Clone()
		if !f(node) {
			return
		}
	}
}

//Watch 监控节点变化
func (c *Cluster) Watch() conf.IWatcher {
	w := newCWatcher(c)
	c.watchers.Set(w.id, w)
	w.notify(c.current)
	return w
}

//Close 关闭当前集群管理
func (c *Cluster) Close() error {
	close(c.closeCh)
	return nil
}

//removeWatcher 移除监控器
func (c *Cluster) removeWatcher(id string) {
	c.watchers.Remove(id)
}

//-------------------------------------内部处理-----------------------------------
func (c *Cluster) load() error {
	if err := c.getCluster(); err != nil {
		return err
	}
	errs := make(chan error, 1)
	go func() {
		err := c.watchCluster()
		if err != nil {
			errs <- err
		}
	}()
	select {
	case err := <-errs:
		return err
	case <-time.After(time.Millisecond * 500):
		return nil
	}
}
func (c *Cluster) getCluster() error {
	path := c.GetServerPubPath()
	current := c.current
	children, _, err := c.registry.GetChildren(path)
	if err != nil {
		return err
	}
	sort.Strings(children)

	//移除所有已下线的节点
	c.nodes.RemoveIterCb(func(key string, v interface{}) bool {
		removeNow := true
		for _, name := range children {
			if name == key {
				removeNow = false
				break
			}
		}
		return removeNow
	})

	//设置或添加在线节点
	for i, name := range children {
		node := NewCNode(name, c.GetServerID(), i)
		c.nodes.Set(name, node)
		if node.IsCurrent() {
			current = node
		}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.current = current

	//通知所有订阅者
	c.watchers.IterCb(func(key string, v interface{}) bool {
		watcher := v.(*cwatcher)
		watcher.notify(c.current)
		return true
	})
	return nil
}
func (c *Cluster) watchCluster() error {
	wc, err := watcher.NewChildWatcherByRegistry(c.registry, []string{c.GetServerPubPath()}, logger.New("watch.server"))
	if err != nil {
		return err
	}
	notify, err := wc.Start()
	if err != nil {
		return err
	}
LOOP:
	for {
		select {
		case <-c.closeCh:
			break LOOP
		case <-notify:
			c.getCluster()
		}
	}
	return nil
}
