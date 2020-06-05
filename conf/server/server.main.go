package server

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//MainConf 服务器主配置
type MainConf struct {
	mainConf    *conf.JSONConf
	version     int32
	subConfs    map[string]conf.JSONConf
	currentNode *CNode
	registry    registry.IRegistry
	conf.IPub
	cluster conf.ICluster
	closeCh chan struct{}
}

//NewMainConf 管理服务器的主配置信息
func NewMainConf(platName string, systemName string, serverType string, clusterName string, rgst registry.IRegistry) (s *MainConf, err error) {
	s = &MainConf{
		registry: rgst,
		IPub:     NewPub(platName, systemName, serverType, clusterName),
		subConfs: make(map[string]conf.JSONConf),
		closeCh:  make(chan struct{}),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//load 加载配置
func (c *MainConf) load() (err error) {

	//获取主配置
	conf, err := c.getValue(c.GetMainPath())
	if err != nil {
		return err
	}
	c.mainConf = conf
	c.version = conf.GetVersion()

	//获取子配置
	c.subConfs, err = c.getSubConf(c.GetMainPath())
	if err != nil {
		return err
	}

	//获取节点信息
	c.loadCluster()
	return nil
}

func (c *MainConf) getSubConf(path string) (map[string]conf.JSONConf, error) {
	confs, _, err := c.registry.GetChildren(path)
	if err != nil {
		return nil, err
	}
	values := make(map[string]conf.JSONConf)
	for _, p := range confs {
		currentPath := registry.Join(path, p)
		value, err := c.getValue(currentPath)
		if err != nil {
			return nil, err
		}

		children, err := c.getSubConf(currentPath)
		if err != nil {
			return nil, err
		}
		for k, v := range children {
			values[registry.Join(p, k)] = v
		}
		if len(children) == 0 {
			values[p] = *value
		}
	}
	return values, nil
}

func (c *MainConf) getValue(path string) (*conf.JSONConf, error) {
	data, version, err := c.registry.GetValue(path)
	if err != nil {
		return nil, fmt.Errorf("获取配置出错 %s %w", path, err)
	}

	rdata, err := decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("%s[%s]解密子配置失败:%w", path, data, err)
	}
	if len(rdata) == 0 {
		rdata = []byte("{}")
	}
	childConf, err := conf.NewJSONConf(rdata, version)
	if err != nil {
		err = fmt.Errorf("%s[%s]配置有误:%w", path, data, err)
		return nil, err
	}
	return childConf, nil
}

func (c *MainConf) watchCluster() error {
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
			c.loadCluster()
		}
	}
	return nil
}
func (c *MainConf) loadCluster() error {
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
	return nil
}

//IsTrace 是否跟踪请求或响应
func (c *MainConf) IsTrace() bool {
	return c.mainConf.GetString("trace", "true") == "true"
}

//GetRegistry 获取注册中心
func (c *MainConf) GetRegistry() registry.IRegistry {
	return c.registry
}

//IsStarted 当前服务是否已启动
func (c *MainConf) IsStarted() bool {
	return c.mainConf.GetString("status", "start") == "start"
}

//GetVersion 获取版本号
func (c *MainConf) GetVersion() int32 {
	return c.version
}

//GetCluster 获取集群中的所有节点
func (c *MainConf) GetCluster() conf.ICluster {
	return c.cluster.Clone()
}

//GetMainConf 获取当前主配置
func (c *MainConf) GetMainConf() *conf.JSONConf {
	return c.mainConf
}

//GetMainObject 获取主配置信息
func (c *MainConf) GetMainObject(v interface{}) (int32, error) {
	conf := c.GetMainConf()
	if err := conf.Unmarshal(&v); err != nil {
		err = fmt.Errorf("获取主配置失败:%v", err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//GetClusterNode 获取当前集群节点信息
func (c *MainConf) GetClusterNode() conf.ICNode {
	return c.currentNode
}

//GetSubConf 指定子配置
func (c *MainConf) GetSubConf(name string) (*conf.JSONConf, error) {
	if v, ok := c.subConfs[name]; ok {
		return &v, nil
	}
	return nil, conf.ErrNoSetting
}

//GetSubObject 获取子配置信息
func (c *MainConf) GetSubObject(name string, v interface{}) (int32, error) {
	conf, err := c.GetSubConf(name)
	if err != nil {
		return 0, err
	}

	if err := conf.Unmarshal(&v); err != nil {
		err = fmt.Errorf("获取%s配置失败:%v", name, err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//Has 是否存在子配置
func (c *MainConf) Has(names ...string) bool {
	for _, name := range names {
		_, ok := c.subConfs[name]
		if ok {
			return true
		}
	}
	return false
}

//Iter 迭代所有配置
func (c *MainConf) Iter(f func(path string, conf *conf.JSONConf) bool) {
	for path, v := range c.subConfs {
		if !f(path, &v) {
			break
		}
	}
}

//Close 关闭清理资源
func (c *MainConf) Close() error {
	close(c.closeCh)
	return nil
}
