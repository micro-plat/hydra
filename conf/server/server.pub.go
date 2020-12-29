package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
)

var _ conf.IServerPub = &ServerPub{}

//ServerPub 系统发布路径
type ServerPub struct {
	platName    string
	sysName     string
	serverType  string
	clusterName string
	serverID    string
}

//NewServerPub 构建服务发布路径信息
func NewServerPub(platName string, sysName string, serverType string, clusterName string) *ServerPub {
	return &ServerPub{
		platName:    platName,
		sysName:     sysName,
		serverType:  serverType,
		clusterName: clusterName,
		serverID:    global.GetMatchineCode() + strconv.Itoa(os.Getpid()),
	}
}

//GetServerRoot 获取服务器根路径
func (c *ServerPub) GetServerRoot() string {
	return registry.Join(c.platName, c.sysName, c.serverType)
}

//GetServerPath 获取配置路径
func (c *ServerPub) GetServerPath() string {
	return registry.Join(c.platName, c.sysName, c.serverType, c.clusterName, "conf")
}

//GetSubConfPath 获取子配置路径
func (c *ServerPub) GetSubConfPath(name ...string) string {
	l := []string{c.GetServerPath()}
	l = append(l, name...)
	return registry.Join(l...)
}

//GetRPCServicePubPath 获取服务发布跟路径
func (c *ServerPub) GetRPCServicePubPath(svName string) string {
	return registry.Join(c.platName, "services", c.serverType, svName, "providers")
}

//GetServicePubPath 获取服务发布跟路径
func (c *ServerPub) GetServicePubPath() string {
	return registry.Join(c.platName, "services", c.serverType, "providers")
}

//GetDNSPubPath 获取DNS服务路径
func (c *ServerPub) GetDNSPubPath(svName string) string {
	return registry.Join(global.Def.GetDNSRoot(), svName)
}

//GetServerPubPath 获取服务器发布的跟路径
func (c *ServerPub) GetServerPubPath(clustName ...string) string {
	if len(clustName) == 0 {
		return registry.Join(c.platName, c.sysName, c.serverType, c.clusterName, "servers")
	}
	return registry.Join(c.platName, c.sysName, c.serverType, clustName[0], "servers")
}

//GetServerID 获取当前服务的集群编号
func (c *ServerPub) GetServerID() string {
	return c.serverID
}

//GetPlatName 获取平台名称
func (c *ServerPub) GetPlatName() string {
	return c.platName
}

//GetSysName 获取系统名称
func (c *ServerPub) GetSysName() string {
	return c.sysName
}

//GetServerType 获取服务器类型
func (c *ServerPub) GetServerType() string {
	return c.serverType
}

//GetClusterName 获取集群名称
func (c *ServerPub) GetClusterName() string {
	return c.clusterName
}

//GetServerName 获取服务器名称
func (c *ServerPub) GetServerName() string {
	return fmt.Sprintf("%s(%s)", c.sysName, c.clusterName)
}

//AllowGray 是否允许灰度到其它集群
func (c *ServerPub) AllowGray() bool {
	return c.serverType == global.API ||
		c.serverType == global.Web
}
