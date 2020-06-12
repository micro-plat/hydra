package server

import (
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

type Cluster struct {
	currentNode *CNode
	cluster     []*CNode
	registry    registry.IRegistry
	closeCh     chan struct{}
	conf.IPub
}

//NewCluster 管理服务器的主配置信息
func NewCluster(platName string, systemName string, serverType string, clusterName string, rgst registry.IRegistry) (s *MainConf, err error) {
	s = &MainConf{
		registry: rgst,
		IPub:     NewPub(platName, systemName, serverType, clusterName),
		closeCh:  make(chan struct{}),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}
func (c *Cluster) Iter(f func(conf.ICNode) bool) {

}
func (c *Cluster) Current() conf.ICNode {
	return nil
}
func (c *Cluster) Clone() conf.ICluster {
	return nil
}

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
	cnodes := make([]*CNode, 0, 2)
	path := c.GetServerPubPath()
	children, _, err := c.registry.GetChildren(path)
	if err != nil {
		return err
	}

	for i, name := range children {
		node := NewCNode(name, c.GetClusterID(), i)
		cnodes = append(cnodes, node)
		if node.IsCurrent() {
			c.currentNode = node
		}
	}
	c.cluster = ClusterNodes(cnodes)
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
