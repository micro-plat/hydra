package config

import (
	"fmt"
	"sync"

	"github.com/micro-plat/hydra/registry"
)

//VarConf 变量信息
type VarConf struct {
	*JSONConf
	varConfPath  string
	varVersion   int32
	varNodeConfs map[string]JSONConf
	registry     registry.IRegistry
	varLock      sync.RWMutex
}

//NewVarConf 构建服务器配置缓存
func NewVarConf(varConfPath string, rgst registry.IRegistry) (s *VarConf, err error) {
	s = &VarConf{
		varConfPath:  varConfPath,
		registry:     rgst,
		varNodeConfs: make(map[string]JSONConf),
	}
	if err = s.loadVarNodeConf(); err != nil {
		return
	}
	return s, nil
}

//初始化子节点配置
func (c *VarConf) loadVarNodeConf() (err error) {

	//检查跟路径是否存在
	if b, err := c.registry.Exists(c.varConfPath); err == nil && !b {
		return nil
	}

	//获取第一级目录
	var varfirstNodes []string
	varfirstNodes, c.varVersion, err = c.registry.GetChildren(c.varConfPath)
	if err != nil {
		return err
	}

	for _, p := range varfirstNodes {
		//获取第二级目录
		firstNodePath := registry.Join(c.varConfPath, p)
		varSecondChildren, _, err := c.registry.GetChildren(firstNodePath)
		if err != nil {
			return err
		}

		//获取二级目录的值
		for _, node := range varSecondChildren {
			nodePath := registry.Join(firstNodePath, node)
			data, version, err := c.registry.GetValue(nodePath)
			if err != nil {
				return err
			}
			rdata, err := decrypt(data)
			if err != nil {
				return err
			}
			varConf, err := NewJSONConf(rdata, version)
			if err != nil {
				err = fmt.Errorf("%s配置有误:%v", nodePath, err)
				return err
			}
			c.varNodeConfs[registry.Join(p, node)] = *varConf
		}
	}
	return nil
}

//GetVarVersion 获取var路径版本号
func (c *VarConf) GetVarVersion() int32 {
	return c.varVersion
}

//IterVarConf 迭代所有子配置
func (c *VarConf) IterVarConf(f func(k string, conf *JSONConf) bool) {
	for k, v := range c.varNodeConfs {
		if !f(k, &v) {
			break
		}
	}
}

//GetVarConf 指定配置文件名称，获取var配置信息
func (c *VarConf) GetVarConf(tp string, name string) (*JSONConf, error) {
	c.varLock.RLock()
	defer c.varLock.RUnlock()
	if v, ok := c.varNodeConfs[registry.Join(tp, name)]; ok {
		return &v, nil
	}
	return nil, ErrNoSetting
}

//GetVarConfClone 获取var配置拷贝
func (c *VarConf) GetVarConfClone() map[string]JSONConf {
	c.varLock.RLock()
	defer c.varLock.RUnlock()
	data := make(map[string]JSONConf)
	for k, v := range c.varNodeConfs {
		data[k] = v
	}
	return data
}

//SetVarConf 获取var配置参数
func (c *VarConf) SetVarConf(data map[string]JSONConf) {
	c.varLock.Lock()
	defer c.varLock.Unlock()
	c.varNodeConfs = data
}

//HasVarConf 是否存在子级配置
func (c *VarConf) HasVarConf(tp string, name string) bool {
	_, ok := c.varNodeConfs[registry.Join(tp, name)]
	return ok
}
