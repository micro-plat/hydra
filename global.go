package hydra

import (
	"fmt"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

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

//C 基础组件
var C = components.Def

//OnReady 系统准备好后执行
var OnReady = global.OnReady

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

//ICli 终端命令参数
type ICli = global.ICli
