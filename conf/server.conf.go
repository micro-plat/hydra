package conf

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

type ISystemConf interface {
	GetPlatName() string
	GetSysName() string
	GetServerType() string
	GetClusterName() string
	GetServerName() string
	GetAppConf(v interface{}) error
	Get(key string) interface{}
	Set(key string, value interface{})
}
type IMainConf interface {
	IConf
	GetSystemRootfPath() string
	GetMainConfPath() string
	GetServicePubRootPath(name string) string
	GetServerPubRootPath() string
	IsStop() bool
	ForceRestart() bool
	GetSubObject(name string, v interface{}) (int32, error)
	GetSubConf(name string) (*JSONConf, error)
	HasSubConf(name ...string) bool
	GetSubConfClone() map[string]JSONConf
	SetSubConf(data map[string]JSONConf)
	IterSubConf(f func(k string, conf *JSONConf) bool)
}
type IVarConf interface {
	GetVarVersion() int32
	GetVarConf(tp string, name string) (*JSONConf, error)
	GetVarObject(tp string, name string, v interface{}) (int32, error)
	HasVarConf(tp string, name string) bool
	GetVarConfClone() map[string]JSONConf
	SetVarConf(map[string]JSONConf)
	IterVarConf(f func(k string, conf *JSONConf) bool)
}

//IServerConf 服务器配置
type IServerConf interface {
	ISystemConf
	IMainConf
	IVarConf
}

//ServerConf 服务器配置信息
type ServerConf struct {
	*JSONConf
	platName     string
	sysName      string
	serverType   string
	clusterName  string
	mainConfpath string
	varConfPath  string
	*metadata
	varVersion   int32
	subNodeConfs map[string]JSONConf
	varNodeConfs map[string]JSONConf
	registry     registry.IRegistry
	varLock      sync.RWMutex
	subLock      sync.RWMutex
}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, mainConfRaw []byte, mainConfVersion int32, rgst registry.IRegistry) (s *ServerConf, err error) {

	sections := strings.Split(strings.Trim(mainConfpath, rgst.GetSeparator()), rgst.GetSeparator())
	if len(sections) != 5 {
		err = fmt.Errorf("conf配置文件格式错误，格式:/platName/sysName/serverType/clusterName/conf 当前值：%s", mainConfpath)
		return
	}
	s = &ServerConf{
		metadata:     &metadata{},
		mainConfpath: mainConfpath,
		platName:     sections[0],
		sysName:      sections[1],
		serverType:   sections[2],
		clusterName:  sections[3],
		varConfPath:  registry.Join("/", sections[0], "var"),
		registry:     rgst,
		subNodeConfs: make(map[string]JSONConf),
		varNodeConfs: make(map[string]JSONConf),
	}
	rdata, err := decrypt(mainConfRaw)
	if err != nil {
		return nil, err
	}
	//初始化主配置
	if s.JSONConf, err = NewJSONConf(rdata, mainConfVersion); err != nil {
		err = fmt.Errorf("%s配置有误:%v", mainConfpath, err)
		return nil, err
	}
	if s.GetString("status", "start") != "start" && s.GetString("status", "start") != "stop" && s.GetString("status", "start") != "restart" {
		err = fmt.Errorf("%s配置有误:status的值只能是'start','stop'或 'restart'", mainConfpath)
		return nil, err
	}
	if err = s.loadChildNodeConf(); err != nil {
		return
	}
	if err = s.loadVarNodeConf(); err != nil {
		return
	}
	return s, nil
}

//初始化子节点配置
func (c *ServerConf) loadChildNodeConf() error {
	paths, _, err := c.registry.GetChildren(c.mainConfpath)
	if err != nil {
		return err
	}
	for _, p := range paths {
		childConfPath := registry.Join(c.mainConfpath, p)
		data, version, err := c.registry.GetValue(childConfPath)
		if err != nil {
			return err
		}
		rdata, err := decrypt(data)
		if err != nil {
			return err
		}
		childConf, err := NewJSONConf(rdata, version)
		if err != nil {
			err = fmt.Errorf("%s配置有误:%v", childConfPath, err)
			return err
		}
		c.subNodeConfs[p] = *childConf
	}
	return nil
}

//初始化子节点配置
func (c *ServerConf) loadVarNodeConf() (err error) {

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

//IsStop 当前服务是否已停止
func (c *ServerConf) IsStop() bool {
	return c.GetString("status", "start") != "start" && c.GetString("status", "start") != "restart"
}

//ForceRestart 强制重启
func (c *ServerConf) ForceRestart() bool {
	return c.GetString("status", "start") == "restart"
}

//GetMainConfPath 获取主配置文件路径
func (c *ServerConf) GetMainConfPath() string {
	return registry.Join("/", c.mainConfpath)
}

//GetSystemRootfPath 获取系统根路径
func (c *ServerConf) GetSystemRootfPath() string {
	return registry.Join("/", c.platName, c.sysName, c.serverType, c.clusterName)
}

//GetServicePubRootPath 获取服务发布跟路径
func (c *ServerConf) GetServicePubRootPath(svName string) string {
	return registry.Join("/", c.platName, "services", c.serverType, svName, "providers")
}

//GetServerPubRootPath 获取服务器发布的跟路径
func (c *ServerConf) GetServerPubRootPath() string {
	return registry.Join("/", c.GetSystemRootfPath(), "servers")
}

//GetAppConf 获取系统配置
func (c *ServerConf) GetAppConf(v interface{}) error {
	_, err := c.GetSubObject("app", v)
	if err != nil {
		return fmt.Errorf("获取app配置出错:%v", err)
	}
	return err
}

//GetSubObject 获取子系统配置
func (c *ServerConf) GetSubObject(name string, v interface{}) (int32, error) {
	conf, err := c.GetSubConf(name)
	if err != nil {
		return 0, err
	}
	if err := conf.Unmarshal(&v); err != nil {
		err = fmt.Errorf("获取%s配置失败:%v", name, err)
		return 0, err
	}

	return conf.version, nil
}

//GetVarVersion 获取var路径版本号
func (c *ServerConf) GetVarVersion() int32 {
	return c.varVersion
}

//GetSubConf 指定配置文件名称，获取系统配置信息
func (c *ServerConf) GetSubConf(name string) (*JSONConf, error) {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	if v, ok := c.subNodeConfs[name]; ok {
		return &v, nil
	}
	return nil, ErrNoSetting
}

//GetSubConfClone 获取sub配置拷贝
func (c *ServerConf) GetSubConfClone() map[string]JSONConf {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	data := make(map[string]JSONConf)
	for k, v := range c.subNodeConfs {
		data[k] = v
	}
	return data
}

//SetSubConf 获取sub配置参数
func (c *ServerConf) SetSubConf(data map[string]JSONConf) {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	c.subNodeConfs = data
}

//HasSubConf 是否存在子级配置
func (c *ServerConf) HasSubConf(names ...string) bool {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	for _, name := range names {
		_, ok := c.subNodeConfs[name]
		if ok {
			return true
		}
	}
	return false
}

//IterSubConf 迭代所有子配置
func (c *ServerConf) IterSubConf(f func(k string, conf *JSONConf) bool) {
	for k, v := range c.subNodeConfs {
		if !f(k, &v) {
			break
		}
	}
}

//IterVarConf 迭代所有子配置
func (c *ServerConf) IterVarConf(f func(k string, conf *JSONConf) bool) {
	for k, v := range c.varNodeConfs {
		if !f(k, &v) {
			break
		}
	}
}

//GetVarConf 指定配置文件名称，获取var配置信息
func (c *ServerConf) GetVarConf(tp string, name string) (*JSONConf, error) {
	c.varLock.RLock()
	defer c.varLock.RUnlock()
	if v, ok := c.varNodeConfs[registry.Join(tp, name)]; ok {
		return &v, nil
	}
	return nil, ErrNoSetting
}

//GetVarConfClone 获取var配置拷贝
func (c *ServerConf) GetVarConfClone() map[string]JSONConf {
	c.varLock.RLock()
	defer c.varLock.RUnlock()
	data := make(map[string]JSONConf)
	for k, v := range c.varNodeConfs {
		data[k] = v
	}
	return data
}

//SetVarConf 获取var配置参数
func (c *ServerConf) SetVarConf(data map[string]JSONConf) {
	c.varLock.Lock()
	defer c.varLock.Unlock()
	c.varNodeConfs = data
}

//GetVarObject 指定配置文件名称，获取var配置信息
func (c *ServerConf) GetVarObject(tp string, name string, v interface{}) (int32, error) {
	conf, err := c.GetVarConf(tp, name)
	if err != nil {
		return 0, err
	}
	if err := conf.Unmarshal(&v); err != nil {
		err = fmt.Errorf("获取/%s/%s配置失败:%v", tp, name, err)
		return 0, err
	}
	return conf.version, nil
}

//HasVarConf 是否存在子级配置
func (c *ServerConf) HasVarConf(tp string, name string) bool {
	_, ok := c.varNodeConfs[registry.Join(tp, name)]
	return ok
}

//GetPlatName 获取平台名称
func (c *ServerConf) GetPlatName() string {
	return c.platName
}

//GetSysName 获取系统名称
func (c *ServerConf) GetSysName() string {
	return c.sysName
}

//GetServerType 获取服务器类型
func (c *ServerConf) GetServerType() string {
	return c.serverType
}

//GetClusterName 获取集群名称
func (c *ServerConf) GetClusterName() string {
	return c.clusterName
}

//GetServerName 获取服务器名称
func (c *ServerConf) GetServerName() string {
	return fmt.Sprintf("%s.%s(%s)", c.sysName, c.clusterName, c.serverType)
}
