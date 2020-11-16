package global

import (
	"os"
	"path/filepath"

	"github.com/micro-plat/lib4go/logger"
)

//AppName 当前应用程序的名称
var AppName string = filepath.Base(os.Args[0])

//Version 版本号
var Version string = "1.0.0"

//Usage 用途
var Usage string = filepath.Base(os.Args[0]) + "(A new hydra application)"

//API api服务器
const API = "api"

//Web web服务器
const Web = "web"

//RPC rpc服务器
const RPC = "rpc"

//WS websocket服务器
const WS = "ws"

//CRON cron服务器
const CRON = "cron"

//MQC mqc服务器
const MQC = "mqc"

//ServerTypes 支持的所有服务器类型
var ServerTypes = []string{}

//IGlobal 应用程序信息
type IGlobal interface {

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

	//GetTracePort 当Trace为web时候，需要设置TracePort
	GetTracePort() string

	//Log 获取日志组件
	Log() logger.ILogger

	//HasServerType 是否包含指定的服务器类型
	HasServerType(tp string) bool

	//ClosingNotify 获取系统关闭通知
	ClosingNotify() chan struct{}

	//Close 关闭应用
	Close()
}

//Current 当前应用程序信息
func Current() IGlobal {
	return Def
}
