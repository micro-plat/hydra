package conf

import (
	"errors"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//IServerConf 主配置信息
type IServerConf interface {
	IServerPub
	IsStarted() bool
	IsTrace() bool
	GetRootConf() *RawConf
	GetMainObject(v interface{}) (int32, error)
	GetSubConf(name string) (*RawConf, error)
	GetSubObject(name string, v interface{}) (int32, error)
	GetRegistry() registry.IRegistry
	GetVersion() int32
	GetCluster(clustName ...string) (ICluster, error)
	Has(names ...string) bool
	Iter(f func(path string, conf *RawConf) bool)
	Close() error
}

//ICNode 集群节点配置
type ICNode interface {
	IsAvailable() bool
	GetHost() string
	GetPort() string
	GetServerID() string
	GetName() string
	IsCurrent() bool
	GetIndex() int
	//IsMaster 判断当前节点是否是主服务器
	IsMaster(i int) bool
	Clone() ICNode
}

//ICluster 集群信息
type ICluster interface {
	Iter(f func(ICNode) bool)
	Current() ICNode
	Watch() IWatcher
	Next() (ICNode, bool)
	GetType() string
	Close() error
}

//IServerPub 发布路径服务
type IServerPub interface {
	GetServerPath() string
	GetSubConfPath(name ...string) string
	GetRPCServicePubPath(svName string) string
	GetServicePubPath() string
	GetDNSPubPath(svName string) string
	GetServerPubPath(clustName ...string) string
	GetServerID() string
	GetPlatName() string
	GetSysName() string
	GetServerType() string
	GetClusterName() string
	GetServerName() string
	AllowGray() bool
}

//IVarConf 变量配置
type IVarConf interface {
	IVarPub
	GetVersion() int32
	GetConf(tp string, name string) (*RawConf, error)
	GetConfVersion(tp string, name string) (int32, error)
	GetObject(tp string, name string, v interface{}) (int32, error)
	GetClone() IVarConf
	Has(tp string, name string) bool
	Iter(f func(k string, conf *RawConf) bool)
}

//IPub 发布路径服务
type IVarPub interface {
	GetVarPath(p ...string) string
	GetRLogPath() string
}

//IWatcher 集群监控
type IWatcher interface {
	Notify() chan ICNode
	Close() error
}
