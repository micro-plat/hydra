package application

import (
	"os"
	"path/filepath"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
)

//AppName 当前应用程序的名称
var AppName string = filepath.Base(os.Args[0])

//Version 版本号
var Version string = "1.0.0"

//IApplication 应用程序信息
type IApplication interface {
	GetCMD() string
	GetHandler(tp string, service string) context.IHandler
	Server(tp string) server.IServerConf
	CurrentContext() context.IContext

	//GetRegistryAddr 注册中心
	GetRegistryAddr() string

	//GetPlatName 平台名称
	GetPlatName() string

	//GetSysName 系统名称
	GetSysName() string

	//GetServerTypes 服务器类型
	GetServerTypes() []string

	//GetClusterName 集群名称
	GetClusterName() string

	//GetTrace 显示请求与响应信息
	GetTrace() string
}

//Current 当前应用程序信息
func Current() IApplication {
	return DefApp
}
