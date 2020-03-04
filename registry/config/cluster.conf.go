package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

//ClusterConf 集群配置
type ClusterConf struct {
	*JSONConf
	platName     string
	sysName      string
	serverType   string
	clusterName  string
	clusterID    string
	mainConfpath string
	*metadata
	subNodeConfs map[string]JSONConf
	registry     registry.IRegistry
	subLock      sync.RWMutex
}

//NewClusterConf 指定集群节点的配置路径，注册中心对象构建集群配置缓存
func NewClusterConf(mainConfpath string, rgst registry.IRegistry) (s *ClusterConf, err error) {
	sections := strings.Split(strings.Trim(mainConfpath, rgst.GetSeparator()), rgst.GetSeparator())
	if len(sections) != 5 {
		err = fmt.Errorf("conf配置文件格式错误，格式:/platName/sysName/serverType/clusterName/conf 当前值：%s", mainConfpath)
		return
	}
	s = &ClusterConf{
		metadata:     &metadata{},
		mainConfpath: mainConfpath,
		platName:     sections[0],
		sysName:      sections[1],
		serverType:   sections[2],
		clusterName:  sections[3],
		clusterID:    utility.GetGUID()[0:8],
		registry:     rgst,
		subNodeConfs: make(map[string]JSONConf),
	}
	mainConfRaw, mainConfVersion, err := rgst.GetValue(mainConfpath)
	if err != nil {
		return nil, err
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
	return s, nil
}

//初始化子节点配置
func (c *ClusterConf) loadChildNodeConf() error {
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

//IsStop 当前服务是否已停止
func (c *ClusterConf) IsStop() bool {
	return c.GetString("status", "start") != "start" && c.GetString("status", "start") != "restart"
}

//GetClusterID 获取当前服务的集群编号
func (c *ClusterConf) GetClusterID() string {
	return c.clusterID
}

//ForceRestart 强制重启
func (c *ClusterConf) ForceRestart() bool {
	return c.GetString("status", "start") == "restart"
}

//GetMainConfPath 获取主配置文件路径
func (c *ClusterConf) GetMainConfPath() string {
	return registry.Join("/", c.mainConfpath)
}

//GetSystemRootfPath 获取系统根路径
func (c *ClusterConf) GetSystemRootfPath(tp ...string) string {
	return registry.Join("/", c.platName, c.sysName,
		types.GetStringByIndex(tp, 0, c.serverType), c.clusterName)
}

//GetServicePubRootPath 获取服务发布跟路径
func (c *ClusterConf) GetServicePubRootPath(svName string) string {
	return registry.Join("/", c.platName, "services", c.serverType, svName, "providers")
}

//GetDNSPubRootPath 获取DNS服务路径
func (c *ClusterConf) GetDNSPubRootPath(svName string) string {
	return registry.Join("/dns", svName)
}

//GetClusterNodes 获取集群其它服务器节点
func (c *ClusterConf) GetClusterNodes(serverType ...string) CNodes {
	cnodes := make([]*CNode, 0, 2)
	path := c.GetServerPubRootPath(serverType...)
	children, _, err := c.registry.GetChildren(path)
	if err != nil {
		return nil
	}

	for i, node := range children {
		cnodes = append(cnodes, NewCNode(path, node, c.clusterID, i))
	}
	return cnodes
}

//GetServerPubRootPath 获取服务器发布的跟路径
func (c *ClusterConf) GetServerPubRootPath(serverType ...string) string {
	return registry.Join("/", c.GetSystemRootfPath(serverType...), "servers")
}

//GetAppConf 获取系统配置
func (c *ClusterConf) GetAppConf(v interface{}) error {
	_, err := c.GetSubObject("app", v)
	if err != nil {
		return fmt.Errorf("获取app配置出错:%v", err)
	}
	return err
}

//GetSubObject 获取子系统配置
func (c *ClusterConf) GetSubObject(name string, v interface{}) (int32, error) {
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

//GetSubConf 指定配置文件名称，获取系统配置信息
func (c *ClusterConf) GetSubConf(name string) (*JSONConf, error) {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	if v, ok := c.subNodeConfs[name]; ok {
		return &v, nil
	}
	return nil, ErrNoSetting
}

//GetSubConfClone 获取sub配置拷贝
func (c *ClusterConf) GetSubConfClone() map[string]JSONConf {
	c.subLock.RLock()
	defer c.subLock.RUnlock()
	data := make(map[string]JSONConf)
	for k, v := range c.subNodeConfs {
		data[k] = v
	}
	return data
}

//SetSubConf 获取sub配置参数
func (c *ClusterConf) SetSubConf(data map[string]JSONConf) {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	c.subNodeConfs = data
}

//HasSubConf 是否存在子级配置
func (c *ClusterConf) HasSubConf(names ...string) bool {
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
func (c *ClusterConf) IterSubConf(f func(k string, conf *JSONConf) bool) {
	for k, v := range c.subNodeConfs {
		if !f(k, &v) {
			break
		}
	}
}

//GetPlatName 获取平台名称
func (c *ClusterConf) GetPlatName() string {
	return c.platName
}

//GetSysName 获取系统名称
func (c *ClusterConf) GetSysName() string {
	return c.sysName
}

//GetServerType 获取服务器类型
func (c *ClusterConf) GetServerType() string {
	return c.serverType
}

//GetClusterName 获取集群名称
func (c *ClusterConf) GetClusterName() string {
	return c.clusterName
}

//GetServerName 获取服务器名称
func (c *ClusterConf) GetServerName() string {
	return fmt.Sprintf("%s.%s(%s)", c.sysName, c.clusterName, c.serverType)
}
