package server

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
)

//MainConf 服务器主配置
type MainConf struct {
	mainConf *conf.JSONConf
	version  int32
	subConfs map[string]conf.JSONConf
	registry registry.IRegistry
	conf.IPub
}

//NewMainConf 管理服务器的主配置信息
func NewMainConf(platName string, systemName string, serverType string, clusterName string, rgst registry.IRegistry) (s *MainConf, err error) {
	s = &MainConf{
		registry: rgst,
		IPub:     NewPub(platName, systemName, serverType, clusterName),
		subConfs: make(map[string]conf.JSONConf),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//load 加载配置
func (c *MainConf) load() (err error) {

	var mainConfRaw []byte
	mainConfRaw, c.version, err = c.registry.GetValue(c.GetMainPath())
	if err != nil {
		return fmt.Errorf("无法获取主配置信息 %s %w", c.GetMainPath(), err)
	}
	rdata, err := decrypt(mainConfRaw)
	if err != nil {
		return fmt.Errorf("%s解密失败:%w", c.GetMainPath(), err)
	}
	//初始化主配置
	if c.mainConf, err = conf.NewJSONConf(rdata, c.version); err != nil {
		err = fmt.Errorf("%s配置有误:%v", c.GetMainPath(), err)
		return err
	}

	err = c.getSubConf(c.GetMainPath())
	if err != nil {
		return err
	}
	return nil
}
func (c *MainConf) getSubConf(path string, n ...string) error {
	confs, _, err := c.registry.GetChildren(path)
	if err != nil {
		return err
	}
	for _, p := range confs {
		childConfPath := registry.Join(path, p)
		data, version, err := c.registry.GetValue(childConfPath)
		if err != nil {
			return fmt.Errorf("获取子配置信息出错 %s[%s] %w", path, p, err)
		}

		rdata, err := decrypt(data)
		if err != nil {
			return fmt.Errorf("%s[%s]解密子配置失败:%w", path, p, err)
		}
		if len(rdata) == 0 {
			rdata = []byte("{}")
		}
		childConf, err := conf.NewJSONConf(rdata, version)
		if err != nil {
			err = fmt.Errorf("%s/%s配置有误:%w", path, p, err)
			return err
		}
		nodePath := registry.Trim(registry.Join(types.GetStringByIndex(n, 0, ""), p))
		c.subConfs[nodePath] = *childConf

		if err := c.getSubConf(childConfPath, p); err != nil {
			return err
		}
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

//GetClusterNodes 获取集群中的所有节点
func (c *MainConf) GetClusterNodes() []conf.ICNode {
	cnodes := make([]conf.ICNode, 0, 2)
	path := c.GetServerPubPath()
	children, _, err := c.registry.GetChildren(path)
	if err != nil {
		return nil
	}

	for i, name := range children {
		cnodes = append(cnodes, NewCNode(name, c.GetClusterID(), i))
	}
	return cnodes
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
