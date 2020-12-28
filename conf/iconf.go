package conf

import (
	"errors"

	"github.com/micro-plat/hydra/registry"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//ByInstall 通过安装设置
const ByInstall = "###"

//ByInstallI 通过安装设置
const ByInstallI = -2<<31 - 1

//IServerConf 主配置信息
type IServerConf interface {
	IServerPub

	//IsStarted 服务器是否已启动
	IsStarted() bool

	//IsTrace 是否启动了请求跟踪
	IsTrace() bool

	//GetMainConf 获取服务器原始配置
	GetMainConf() *RawConf

	//GetMainObject 获取服务器原始配置并转换为指定对象
	GetMainObject(out interface{}) (int32, error)

	//GetSubConf 获取子节点原始配置
	GetSubConf(name string) (*RawConf, error)

	//GetSubObject 获取子节点原始配置并转换为指定对象
	GetSubObject(name string, out interface{}) (int32, error)

	//GetRegistry 获取注册中心
	GetRegistry() registry.IRegistry

	//GetVersion 获服主配置版本号
	GetVersion() int32

	//获取所有当前服务器下所有集群名称
	GetClusterNames() []string

	//GetCluster 获取集群信息
	GetCluster(clustName ...string) (ICluster, error)

	//Has 是否包含子配置节点
	Has(names ...string) bool

	//Iter 迭代所有子配置
	Iter(f func(path string, conf *RawConf) bool)

	//Close 关闭当前配置并清理资源
	Close() error
}

//ICNode 集群节点配置
type ICNode interface {
	//IsAvailable 当前节点是否可用
	IsAvailable() bool

	//GetHost 节点提供的服务的服务器名
	GetHost() string

	//GetPort  节点提供的服务的服务端口
	GetPort() string

	//GetNodeID 获取节点编号
	GetNodeID() string

	//GetName 获取节点名称
	GetName() string

	//IsCurrent 是否是当前server对应的节点
	IsCurrent() bool

	//GetIndex 在集群中的索引编号
	GetIndex() int

	//IsMaster 判断当前节点是否是主服务器
	IsMaster(i int) bool

	//Clone 克隆当前节点
	Clone() ICNode
}

//ICluster 集群信息
type ICluster interface {

	//Iter 迭代集群中的所有节点
	Iter(f func(ICNode) bool)

	//Current 获取当前serer对应的节点
	Current() ICNode

	//Watch 监控节群节点变化
	Watch() IWatcher

	//Next 获取下一个节点
	Next() (ICNode, bool)

	//GetServerType 获取集群的服务器类型
	GetServerType() string

	//Len 获取集群节点个数
	Len() int

	//Close 清理当前对象
	Close() error
}

//IServerPub 发布路径服务
type IServerPub interface {
	GetServerRoot() string
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

//IVarPub 发布路径服务
type IVarPub interface {
	GetVarPath(p ...string) string
	GetRLogPath() string
}

//IWatcher 集群监控
type IWatcher interface {
	Notify() chan ICNode
	Close() error
}
