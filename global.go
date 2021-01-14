package hydra

import (
	"fmt"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
)

//C 基础组件
var C = components.Def

//G 全局应用程序配置
var G = global.Def

//S 服务中心
var S services.IService = services.Def

//Conf 配置组件
var Conf creator.IConf = creator.Conf

//CRON CRON服务可进行动态注册管理
var CRON services.ICRON = services.CRON

//MQC MQC服务动态注册管理
var MQC services.IMQC = services.MQC

//IContext 请求上下文
type IContext = context.IContext

//Installer 安装程序
var Installer = global.Installer

//InstallCli  配置处理相关的终端参数
var InstallCli global.ICustomCli = global.InstallCli

//RunCli 配置处理相关的终端参数
var RunCli global.ICustomCli = global.RunCli

//ConfCli 配置处理相关的终端参数
var ConfCli global.ICustomCli = global.ConfCli

//DBCli 配置处理相关的终端参数
var DBCli global.ICustomCli = global.DBCli

//OnReadying 在系统准备好后首先执行
var OnReadying = global.OnReadying

//OnReady 系统准备好后稍后执行
var OnReady = global.OnReady

//IAPPConf 服务器配置信息
type IAPPConf = app.IAPPConf

//ByInstall 通过安装设置
const ByInstall = conf.ByInstall

//ByInstallI 通过安装设置
const ByInstallI = conf.ByInstallI

//FlagOption 配置选项
type FlagOption = global.FlagOption

//WithFlag 添加字符串flag
var WithFlag = global.WithFlag

//WithBoolFlag 设置bool参数
var WithBoolFlag = global.WithBoolFlag

//WithSliceFlag 设置数组参数
var WithSliceFlag = global.WithSliceFlag

//ICli 终端命令参数
type ICli = global.ICli

//MQConf 消息队列名称管理
var MQConf = global.MQConf

//LocalIP 获取本机IP
var LocalIP = global.LocalIP

//IsDebug 是否是debug模式
var IsDebug = global.Def.IsDebug

//Server 通过服务类型从全局缓存中获取服务配置
func Server(tp string) app.IAPPConf {
	s, err := app.Cache.GetAPPConf(tp)
	if err == nil {
		return s
	}
	panic(fmt.Errorf("[%s]服务器未启动:%w", tp, err))
}

//CurrentContext 获取当前请求上下文
func CurrentContext() context.IContext {
	return context.Current()
}

func init() {
	OnReady(func() error {
		if !registry.Support(global.Def.RegistryAddr) {
			return fmt.Errorf("不支持%s作为注册中心", global.Def.RegistryAddr)
		}
		return nil
	})
}
