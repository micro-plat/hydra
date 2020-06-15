package server

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//MainConf 服务器主配置
type MainConf struct {
	mainConf *conf.JSONConf
	version  int32
	subConfs map[string]conf.JSONConf
	registry registry.IRegistry
	conf.IPub
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
	conf, err := getValue(c.registry, c.GetMainPath())
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
		value, err := getValue(c.registry, currentPath)
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

//GetCluster 获取集群信息
func (c *MainConf) GetCluster() conf.ICluster {
	cluster, err := getCluster(c.IPub, c.registry)
	if err != nil {
		panic(err)
	}
	return cluster
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
func getValue(registry registry.IRegistry, path string) (*conf.JSONConf, error) {
	data, version, err := registry.GetValue(path)
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
