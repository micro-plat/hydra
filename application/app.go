package application

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
)

//DefApp 默认app
var DefApp = &application{}

type application struct {

	//registryAddr 集群地址
	RegistryAddr string

	//PlatName 平台名称
	PlatName string

	//SysName 系统名称
	SysName string

	//ServerTypes 服务器类型
	ServerTypes []string

	//ServerTypeNames 服务类型名称
	ServerTypeNames string

	//ClusterName 集群名称
	ClusterName string

	//Name 服务器请求名称
	Name string

	//Trace 显示请求与响应信息
	Trace string
}

func (a *application) Bind() error {
	return nil
}

func (a *application) GetCMD() string {
	return ""
}

//GetHandler 获取服务对应的处理函数
func (a *application) GetHandler(tp string, service string) context.IHandler {
	return nil
}

//Server 获取服务器配置信息
func (a *application) Server(tp string) server.IServerConf {
	s, err := server.Cache.GetServerConf(tp)
	if err == nil {
		return s
	}
	panic(fmt.Errorf("[%s]服务器未启动:%w", tp, err))
}

//CurrentContext 获取当前请求上下文
func (a *application) CurrentContext() context.IContext {
	return nil
}

//GetRegistryAddr 注册中心
func (a *application) GetRegistryAddr() string {
	return a.RegistryAddr
}

//GetPlatName 平台名称
func (a *application) GetPlatName() string {
	return a.PlatName
}

//GetSysName 系统名称
func (a *application) GetSysName() string {
	return a.SysName
}

//GetServerTypes 服务器类型
func (a *application) GetServerTypes() []string {
	return a.ServerTypes
}

//GetClusterName 集群名称
func (a *application) GetClusterName() string {
	return a.ClusterName
}

//GetTrace 显示请求与响应信息
func (a *application) GetTrace() string {
	return a.Trace
}
