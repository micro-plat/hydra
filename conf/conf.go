package conf

import (
	"errors"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//IMainConf 主配置信息
type IMainConf interface {
	IPub
	IsStarted() bool
	IsTrace() bool
	GetMainConf() *JSONConf
	GetMainObject(v interface{}) (int32, error)
	GetSubConf(name string) (*JSONConf, error)
	GetSubObject(name string, v interface{}) (int32, error)
	GetRegistry() registry.IRegistry
	GetVersion() int32
	GetCluster() ICluster
	Has(names ...string) bool
	Iter(f func(path string, conf *JSONConf) bool)
	Close() error
}

//ICNode 集群节点配置
type ICNode interface {
	IsAvailable() bool
	GetHost() string
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
}

//IPub 发布路径服务
type IPub interface {
	GetMainPath() string
	GetSubConfPath(name ...string) string
	GetServicePubPathByService(svName string) string
	GetServicePubPath() string
	GetDNSPubPath(svName string) string
	GetServerPubPath() string
	GetServerID() string
	GetPlatName() string
	GetSysName() string
	GetServerType() string
	GetClusterName() string
	GetServerName() string
	GetVarPath(p ...string) string
}

//IVarConf 变量配置
type IVarConf interface {
	GetVersion() int32
	GetConf(tp string, name string) (*JSONConf, error)
	GetConfVersion(tp string, name string) (int32, error)
	GetObject(tp string, name string, v interface{}) (int32, error)
	GetClone() IVarConf
	Has(tp string, name string) bool
	Iter(f func(k string, conf *JSONConf) bool)
}

//IWatcher 集群监控
type IWatcher interface {
	Notify() chan ICNode
	Close() error
}
