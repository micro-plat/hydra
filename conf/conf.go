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
	GetClusterNodes() []ICNode
	GetMainConf() *JSONConf
	GetMainObject(v interface{}) (int32, error)
	GetSubConf(name string) (*JSONConf, error)
	GetSubObject(name string, v interface{}) (int32, error)
	GetRegistry() registry.IRegistry
	GetVersion() int32
	Has(names ...string) bool
	Iter(f func(path string, conf *JSONConf) bool)
}

//ICNode 集群节点配置
type ICNode interface {
	GetHost() string
	GetClusterID() string
	GetName() string
	IsCurrent() bool
	GetIndex() int
}

//IPub 发布路径服务
type IPub interface {
	GetMainPath() string
	GetSubConfPath(name ...string) string
	GetServicePubPathByService(svName string) string
	GetServicePubPath() string
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
