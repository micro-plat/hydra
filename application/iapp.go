package application

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

//AppName 当前应用程序的名称
var AppName string = filepath.Base(os.Args[0])

//Version 版本号
var Version string = "1.0.0"

//Usage 用途
var Usage string = " A new hydra application"

//IApplication 应用程序信息
type IApplication interface {
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
func Current() IApplication {
	return DefApp
}

//CheckPrivileges 检查是否有管理员权限
func CheckPrivileges() error {
	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return nil
			}
			return errRootPrivileges
		}
	}
	return errUnsupportedSystem
}

var errUnsupportedSystem = errors.New("Unsupported system")
var errRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")
