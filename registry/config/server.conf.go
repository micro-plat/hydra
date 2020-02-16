package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

//ErrNoSetting 未配置
var ErrNoSetting = errors.New("未配置")

//IClusterConf 集群节点配置
type IClusterConf interface {
	GetPlatName() string
	GetSysName() string
	GetServerType() string
	GetClusterName() string
	GetServerName() string
	GetAppConf(v interface{}) error
	Get(key string) interface{}
	Set(key string, value interface{})
	IConf
	GetSystemRootfPath(server ...string) string
	GetMainConfPath() string
	GetServicePubRootPath(name string) string
	GetServerPubRootPath(serverType ...string) string
	GetDNSPubRootPath(svName string) string
	GetClusterNodes(serverType ...string) CNodes
	IsStop() bool
	ForceRestart() bool
	GetSubObject(name string, v interface{}) (int32, error)
	GetSubConf(name string) (*JSONConf, error)
	HasSubConf(name ...string) bool
	GetSubConfClone() map[string]JSONConf
	SetSubConf(data map[string]JSONConf)
	IterSubConf(f func(k string, conf *JSONConf) bool)
	GetClusterID() string
}

//IVarConf 变量配置
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
	IClusterConf
	IVarConf
}

//ServerConf 服务器配置信息
type ServerConf struct {
	*ClusterConf
	*VarConf
}

//NewServerConf 构建服务器配置缓存
func NewServerConf(mainConfpath string, mainConfRaw []byte, mainConfVersion int32, rgst registry.IRegistry) (s *ServerConf, err error) {
	s = &ServerConf{}
	s.ClusterConf, err = NewClusterConf(mainConfpath, mainConfRaw, mainConfVersion, rgst)
	if err != nil {
		return nil, err
	}
	s.VarConf, err = NewVarConf(mainConfpath, rgst)
	if err != nil {
		return nil, err
	}
	return s, nil

}
