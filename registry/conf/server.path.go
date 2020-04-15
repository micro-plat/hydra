package config

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry"
	"github.com/qxnw/lib4go/utility"
)

var _ IPub = &Pub{}

//IPub 发布路径服务
type IPub interface {
	GetMainPath() string
	GetServicePubPath(svName string) string
	GetDNSPubPath(svName string) string
	GetServerPubPath() string
	GetClusterID() string
	GetPlatName() string
	GetSysName() string
	GetServerType() string
	GetClusterName() string
	GetServerName() string
	GetVarPath() string
}

//Pub 系统发布路径
type Pub struct {
	platName    string
	sysName     string
	serverType  string
	clusterName string
	clusterID   string
}

//Split 根据主配置获取平台名称、系统名称、服务类型、集群名
func Split(mainConfPath string) (platName string, sysName string, serverType string, clusterName string) {
	sections := strings.Split(strings.Trim(mainConfPath, "/"), "/")
	return sections[0], sections[1], sections[2], sections[3]
}

//NewPub 构建服务发布路径信息
func NewPub(platName string, sysName string, serverType string, clusterName string) *Pub {
	return &Pub{
		platName:    platName,
		sysName:     sysName,
		serverType:  serverType,
		clusterName: clusterName,
		clusterID:   utility.GetGUID()[0:8],
	}
}

//GetMainPath 获取配置路径
func (c *Pub) GetMainPath() string {
	return registry.Join(c.platName, c.sysName, c.serverType, c.clusterName, "conf")
}

//GetVarPath 获取var配置路径
func (c *Pub) GetVarPath() string {
	return registry.Join(c.platName, "var")
}

//GetServicePubPath 获取服务发布跟路径
func (c *Pub) GetServicePubPath(svName string) string {
	return registry.Join("/", c.platName, "services", c.serverType, svName, "providers")
}

//GetDNSPubPath 获取DNS服务路径
func (c *Pub) GetDNSPubPath(svName string) string {
	return registry.Join("/dns", svName)
}

//GetServerPubPath 获取服务器发布的跟路径
func (c *Pub) GetServerPubPath() string {
	return registry.Join("/", c.GetMainPath(), "servers")
}

//GetClusterID 获取当前服务的集群编号
func (c *Pub) GetClusterID() string {
	return c.clusterID
}

//GetPlatName 获取平台名称
func (c *Pub) GetPlatName() string {
	return c.platName
}

//GetSysName 获取系统名称
func (c *Pub) GetSysName() string {
	return c.sysName
}

//GetServerType 获取服务器类型
func (c *Pub) GetServerType() string {
	return c.serverType
}

//GetClusterName 获取集群名称
func (c *Pub) GetClusterName() string {
	return c.clusterName
}

//GetServerName 获取服务器名称
func (c *Pub) GetServerName() string {
	return fmt.Sprintf("%s.%s(%s)", c.sysName, c.clusterName, c.serverType)
}
