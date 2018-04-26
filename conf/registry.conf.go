package conf

import (
	"fmt"

	"github.com/micro-plat/hydra/registry"
)

type IRegistryConf interface {
	GetConf(name string) (*JSONConf, error)
	GetObject(name string, v interface{}) (int32, error)
}

//RegistryConf 基于注册中心的配置管理器
type RegistryConf struct {
	registry  registry.IRegistry
	nodeConfs map[string]JSONConf
}

//NewRegistryConf 构建服务器配置缓存
func NewRegistryConf(rgst registry.IRegistry) (s *RegistryConf) {
	return &RegistryConf{
		registry:  rgst,
		nodeConfs: make(map[string]JSONConf),
	}
}

//GetConf 指定配置文件名称，获取系统配置信息
func (c *RegistryConf) GetConf(name string) (*JSONConf, error) {
	if v, ok := c.nodeConfs[name]; ok {
		return &v, nil
	}

	if b, err := c.registry.Exists(name); err == nil && !b {
		return nil, ErrNoSetting
	}
	data, version, err := c.registry.GetValue(name)
	if err != nil {
		return nil, err
	}
	varConf, err := NewJSONConf(data, version)
	if err != nil {
		return nil, fmt.Errorf("%s配置有误:%v", name, err)
	}
	c.nodeConfs[name] = *varConf
	return varConf, nil
}

//GetObject 获取子系统配置
func (c *RegistryConf) GetObject(name string, v interface{}) (int32, error) {
	conf, err := c.GetConf(name)
	if err != nil {
		return 0, err
	}
	if err := conf.Unmarshal(&v); err != nil {
		err = fmt.Errorf("获取%s配置失败:%v", name, err)
		return 0, err
	}
	return conf.version, nil
}
