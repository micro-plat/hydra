package server

import (
	"fmt"
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

func getCluster(pub conf.IPub, rgst registry.IRegistry, clusterName ...string) (s *Cluster, err error) {
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
	conf.IPub
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
func NewCluster(pub conf.IPub, rgst registry.IRegistry, clusterName ...string) (s *Cluster, err error) {
	s = &Cluster{
		IPub:        pub,
		registry:    rgst,
		nodes:       cmap.New(4),
		watchers:    cmap.New(2),
		clusterName: clusterName,
		keyCache:    make([]string, 0, 1),
		closeCh:     make(chan struct{}),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//Current 获取当前节点
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
	nid := atomic.AddInt64(&c.index, 1)
	c.lock.RLock()
	defer c.lock.RUnlock()
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

//Watch 监控节点变化
func (c *Cluster) Watch() conf.IWatcher {
	w := newCWatcher(c)
	c.watchers.Set(w.id, w)
	w.notify(c.current)
	return w
}

//GetType 获取集群类型
func (c *Cluster) GetType() string {
	return c.GetServerType()
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
	path := c.GetServerPubPath(c.clusterName...)
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
		//移除缓存key
		if removeNow {
			c.removeKey(key)
		}
		return removeNow
	})

	//设置或添加在线节点
	for i, name := range children {
		node := NewCNode(name, c.GetServerID(), i)
		if ok, _ := c.nodes.SetIfAbsent(name, node); ok {
			c.addKey(name) //添加到缓存keys中
		}
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
	fmt.Println("11111111111:", []string{c.GetServerPubPath(c.clusterName...)})
	wc, err := watcher.NewChildWatcherByRegistry(c.registry, []string{c.GetServerPubPath(c.clusterName...)}, logger.New("watch.server"))
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
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-c.closeCh:
			break LOOP
		case <-notify:
			c.getCluster()
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
