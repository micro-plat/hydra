package server

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
)

var cluster = cmap.New(2)

func getCluster(pub conf.IServerPub, rgst registry.IRegistry, clusterName ...string) (s *Cluster, err error) {
	_, c, err := cluster.SetIfAbsentCb(pub.GetServerPubPath(clusterName...), func(...interface{}) (interface{}, error) {
		return NewCluster(pub, rgst, clusterName...)
	})
	if err != nil {
		return nil, err
	}
	return c.(*Cluster), nil
}

//Cluster 集群
type Cluster struct {
	conf.IServerPub
	index       int64
	registry    registry.IRegistry
	current     conf.ICNode
	nodes       cmap.ConcurrentMap
	keyCache    []string
	closeCh     chan struct{}
	lock        sync.RWMutex
	watchers    cmap.ConcurrentMap
	clusterName []string
}

//NewCluster 管理服务器的主配置信息
func NewCluster(pub conf.IServerPub, rgst registry.IRegistry, clusterName ...string) (s *Cluster, err error) {
	s = &Cluster{
		IServerPub:  pub,
		registry:    rgst,
		nodes:       cmap.New(4),
		watchers:    cmap.New(2),
		clusterName: clusterName,
		keyCache:    make([]string, 0, 1),
		closeCh:     make(chan struct{}),
	}
	if err = s.getAndWatch(); err != nil {
		return
	}
	return s, nil
}

//Current 获取当前服务器的集群节点
func (c *Cluster) Current() conf.ICNode {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.current == nil {
		return &CNode{}
	}
	return c.current.Clone()
}

//Next 采用轮循的方式获得下一个节点
func (c *Cluster) Next() (conf.ICNode, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if len(c.keyCache) == 0 {
		return nil, false
	}
	nid := atomic.AddInt64(&c.index, 1)
	key := c.keyCache[int(nid%int64(len(c.keyCache)))]

	v, ok := c.nodes.Get(key)
	if ok {
		return v.(conf.ICNode), true
	}
	return nil, false
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

//Len 获取集群节点个数
func (c *Cluster) Len() int {
	return len(c.nodes.Items())
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
func (c *Cluster) getAndWatch() error {
	errs := make(chan error, 1)
	go func() {
		if err := c.watchCluster(); err != nil {
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

	//重新拉取所有节点
	current := c.current
	path := c.GetServerPubPath(c.clusterName...)
	children, _, err := c.registry.GetChildren(path)
	if err != nil {
		return err
	}

	//对节点进行排序
	sort.Strings(children)

	//移除所有已下线的节点
	c.nodes.RemoveIterCb(func(key string, v interface{}) bool {
		//判断当前节点是否需要移除
		removeNow := true
		for _, name := range children {
			if name == key {
				removeNow = false
				break
			}
		}

		//移除已下线节点
		if removeNow {
			node := v.(*CNode)
			if node.IsCurrent() {
				current = &CNode{}
			}
			c.removeKey(key)
		}

		return removeNow
	})

	//修改或添加在线节点
	for i, name := range children {
		node := NewCNode(name, c.GetServerID(), i)
		if ok, _ := c.nodes.SetIfAbsent(name, node); ok {
			c.addKey(name) //添加到缓存keys中
		}
		if node.IsCurrent() {
			current = node
		}
	}

	//更新当前节点
	c.lock.Lock()
	c.current = current
	defer c.lock.Unlock()

	//每次节点变化都通知所有订阅者
	c.watchers.IterCb(func(key string, v interface{}) bool {
		watcher := v.(*cwatcher)
		watcher.notify(c.current)
		return true
	})
	return nil
}

//watchCluster 监听集群变化
func (c *Cluster) watchCluster() error {
	wc, err := watcher.NewChildWatcherByRegistry(c.registry, []string{c.GetServerPubPath(c.clusterName...)}, logger.New("watch.server"))
	if err != nil {
		return err
	}
	notify, err := wc.Start()
	if err != nil {
		return err
	}
	interval := time.Second * 5
	loopPull := time.NewTicker(interval)
	loopPull.Stop()
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-c.closeCh:
			break LOOP
		case <-notify:
			c.getCluster()
			loopPull.Reset(interval)
		case <-loopPull.C:
			c.getCluster()
			loopPull.Stop()
		}
	}
	return nil
}

//addKey 将key添加到顺序缓存表
func (c *Cluster) addKey(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.keyCache = append(c.keyCache, name) //添加到缓存中
}

//removeKey 将key从顺序缓存表中移除
func (c *Cluster) removeKey(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i, k := range c.keyCache {
		if k == name {
			c.keyCache = append(c.keyCache[:i], c.keyCache[i+1:]...)
		}
	}
}
func init() {
	global.Def.AddCloser(func() {
		cluster.RemoveIterCb(func(k string, v interface{}) bool {
			c := v.(*Cluster)
			c.Close()
			return true
		})
	})
}
