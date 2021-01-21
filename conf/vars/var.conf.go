package vars

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var _ conf.IVarConf = &VarConf{}

//EmptyVarConf 空的EmptyVarConf
var EmptyVarConf = &VarConf{
	IVarPub:      nil,
	varConfPath:  "",
	registry:     nil,
	varNodeConfs: make(map[string]conf.RawConf),
}

type cacheObj struct {
	obj     interface{}
	version int32
}

//VarConf 变量信息
type VarConf struct {
	cache cmap.ConcurrentMap
	conf.IVarPub
	varConfPath  string
	varVersion   int32
	varNodeConfs map[string]conf.RawConf
	registry     registry.IRegistry
}

//NewVarConf 构建服务器配置缓存
func NewVarConf(platName string, rgst registry.IRegistry) (s *VarConf, err error) {
	s = &VarConf{
		cache:        cmap.New(3),
		IVarPub:      NewVarPub(platName),
		registry:     rgst,
		varNodeConfs: make(map[string]conf.RawConf),
	}
	s.varConfPath = s.GetVarPath()
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//load 加载所有配置项
func (c *VarConf) load() (err error) {
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
			rdata, err := conf.Decrypt(data)
			if err != nil {
				return err
			}
			varConf, err := conf.NewByText(rdata, version)
			if err != nil {
				err = fmt.Errorf("%s配置有误:%v", nodePath, err)
				return err
			}
			c.varNodeConfs[registry.Join(p, node)] = *varConf
		}
	}
	return nil
}

//GetVersion 获取数据版本号
func (c *VarConf) GetVersion() int32 {
	return c.varVersion
}

//GetConf 指定配置文件名称，获取var配置信息
func (c *VarConf) GetConf(tp string, name string) (*conf.RawConf, error) {

	if v, ok := c.varNodeConfs[registry.Join(tp, name)]; ok {
		return &v, nil
	}
	return conf.EmptyRawConf, fmt.Errorf("%s %w", registry.Join(tp, name), conf.ErrNoSetting)
}

//GetConfVersion 获取配置的版本号
func (c *VarConf) GetConfVersion(tp string, name string) (int32, error) {
	if v, ok := c.varNodeConfs[registry.Join(tp, name)]; ok {
		return v.GetVersion(), nil
	}
	return 0, fmt.Errorf("%s %w", registry.Join(tp, name), conf.ErrNoSetting)
}

//GetObject 获取子配置信息
func (c *VarConf) GetObject(tp string, name string, v interface{}) (int32, error) {
	conf, err := c.GetConf(tp, name)
	if err != nil {
		return 0, err
	}

	if err := conf.ToStruct(&v); err != nil {
		err = fmt.Errorf("获取%s/%s配置失败:%v", tp, name, err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//GetCachedObject　获取已缓存的对象(未缓存时调用GetObject转换对象，并缓存)
func (c *VarConf) GetCachedObject(tp string, name string, v interface{}) (int32, error) {
	key := fmt.Sprintf("%s-%s", tp, name)
	_, obj, err := c.cache.SetIfAbsentCb(key, func(...interface{}) (interface{}, error) {
		ver, err := c.GetObject(tp, name, v)
		if err != nil {
			return nil, err
		}
		ch := &cacheObj{obj: v, version: ver}
		return ch, nil
	})
	if err != nil {
		return 0, err
	}
	ch := obj.(*cacheObj)
	v = ch.obj
	return ch.version, nil
}

//GetClone 获取配置拷贝
func (c *VarConf) GetClone() conf.IVarConf {
	s := &VarConf{
		varVersion:   c.varVersion,
		varConfPath:  c.varConfPath,
		registry:     c.registry,
		varNodeConfs: make(map[string]conf.RawConf),
	}
	for k, v := range c.varNodeConfs {
		s.varNodeConfs[k] = v
	}
	return s
}

//Has 是否存在配置项
func (c *VarConf) Has(tp string, name string) bool {
	_, ok := c.varNodeConfs[registry.Join(tp, name)]
	return ok
}

//Iter 迭代所有子配置
func (c *VarConf) Iter(f func(path string, conf *conf.RawConf) bool) {
	for path, v := range c.varNodeConfs {
		if !f(path, &v) {
			break
		}
	}
}
