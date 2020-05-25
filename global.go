package hydra

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

//Global 全局应用程序配置
var Global = global.DefApp

//Services 服务中心
var Services services.IService = services.DefService

//Conf 配置组件
var Conf creator.IRegistryConf = creator.Conf

//IContext 请求上下文
type IContext = context.IContext

//Component 基础组件
var Component = components.Def
