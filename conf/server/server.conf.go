package server

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
)

//ServerConf 服务器主配置
type ServerConf struct {
	mainConf     *conf.RawConf
	mainVersion  int32
	subConfs     map[string]conf.RawConf
	registry     registry.IRegistry
	clusterNames []string
	conf.IServerPub
	closeCh chan struct{}
}

//NewServerConf 管理服务器的主配置信息
func NewServerConf(platName string, systemName string, serverType string, clusterName string, rgst registry.IRegistry) (s *ServerConf, err error) {
	s = &ServerConf{
		registry:   rgst,
		IServerPub: NewServerPub(platName, systemName, serverType, clusterName),
		subConfs:   make(map[string]conf.RawConf),
		closeCh:    make(chan struct{}),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//IsStarted 当前服务是否已启动
func (c *ServerConf) IsStarted() bool {
	return c.mainConf.GetString("status", "start") == "start"
}

//IsTrace 是否跟踪请求或响应
func (c *ServerConf) IsTrace() bool {
	return c.mainConf.GetBool("trace", false)
}

//GetVersion 获取版本号
func (c *ServerConf) GetVersion() int32 {
	return c.mainVersion
}

//GetMainConf 获取当前主配置
func (c *ServerConf) GetMainConf() *conf.RawConf {
	return c.mainConf
}

//GetMainObject 获取主配置信息
func (c *ServerConf) GetMainObject(v interface{}) (int32, error) {
	conf := c.GetMainConf()
	if err := conf.ToStruct(v); err != nil {
		err = fmt.Errorf("获取主配置失败:%v", err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//GetSubConf 获取子配置
func (c *ServerConf) GetSubConf(name string) (*conf.RawConf, error) {
	if v, ok := c.subConfs[name]; ok {
		return &v, nil
	}

	return nil, conf.ErrNoSetting
}

//GetClusterNames 获取所有当前服务器下所有集群名称
func (c *ServerConf) GetClusterNames() []string {
	return c.clusterNames
}

//GetCluster 获取集群信息
func (c *ServerConf) GetCluster(clustName ...string) (conf.ICluster, error) {
	return getCluster(c.IServerPub, c.registry, clustName...)
}

//GetSubObject 获取子配置信息
func (c *ServerConf) GetSubObject(name string, v interface{}) (int32, error) {
	conf, err := c.GetSubConf(name)
	if err != nil {
		return 0, err
	}

	if err := conf.ToStruct(v); err != nil {
		err = fmt.Errorf("获取%s配置失败:%v", name, err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//Has 是否存在子配置
func (c *ServerConf) Has(names ...string) bool {
	for _, name := range names {
		_, ok := c.subConfs[name]
		if ok {
			return true
		}
	}
	return false
}

//Iter 迭代所有子配置
func (c *ServerConf) Iter(f func(path string, conf *conf.RawConf) bool) {
	for path, v := range c.subConfs {
		if !f(path, &v) {
			break
		}
	}
}

//GetRegistry 获取注册中心
func (c *ServerConf) GetRegistry() registry.IRegistry {
	return c.registry
}

//Close 关闭清理资源
func (c *ServerConf) Close() error {
	close(c.closeCh)
	return nil
}

//load 加载配置
func (c *ServerConf) load() (err error) {

	//获取主配置
	mainPath := c.GetServerPath()
	conf, err := getValue(c.registry, mainPath)
	if err != nil {
		return err
	}
	c.mainConf = conf
	c.mainVersion = conf.GetVersion()

	//获取子配置
	c.subConfs, err = c.getSubConf(c.GetServerPath())
	if err != nil {
		return err
	}

	//获取所有集群名称
	c.clusterNames, _, err = c.registry.GetChildren(c.GetServerRoot())
	if err != nil {
		return err
	}
	return nil
}

func (c *ServerConf) getSubConf(path string) (map[string]conf.RawConf, error) {
	confs, _, err := c.registry.GetChildren(path)
	if err != nil {
		return nil, err
	}

	values := make(map[string]conf.RawConf)
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
func getValue(registry registry.IRegistry, path string) (*conf.RawConf, error) {
	data, version, err := registry.GetValue(path)
	if err != nil {
		return nil, fmt.Errorf("获取配置出错 %s %w", path, err)
	}

	rdata, err := conf.Decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("%s[%s]解密子配置失败:%w", path, data, err)
	}
	if len(rdata) == 0 {
		rdata = []byte("{}")
	}
	childConf, err := conf.NewByText(rdata, version)
	if err != nil {
		err = fmt.Errorf("%s[%s]配置有误:%w", path, data, err)
		return nil, err
	}
	return childConf, nil
}
