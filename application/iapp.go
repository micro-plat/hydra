package application

import (
	"os"
	"path/filepath"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

//AppName 当前应用程序的名称
var AppName string = filepath.Base(os.Args[0])

//Version 版本号
var Version string = "1.0.0"

//RegistryAddr 集群地址
var RegistryAddr string = ""

//PlatName 平台名称
var PlatName string = ""

//SysName 系统名称
var SysName string = ""

//ServerTypes 服务器类型
var ServerTypes []string

//ServerTypeNames 服务类型名称
var ServerTypeNames string

//ClusterName 集群名称
var ClusterName string = ""

//Name 服务器请求名称
var Name string = ""

//Trace 显示请求与响应信息
var Trace string

//Bind 绑定输入参数
func Bind() error {
	//非调试模式时设置日志写协程数为50个
	if !IsDebug {
		logger.AddWriteThread(49)
	}
	return nil
}

//IApplication 应用程序信息
type IApplication interface {
	GetCMD() string
	GetHandler(tp string, service string) context.IHandler
	Server(tp string) server.IServerConf
	CurrentContext() context.IContext
}

//Current 当前应用程序信息
func Current() IApplication {
	return app
}
