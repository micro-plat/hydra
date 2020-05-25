package hydra

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

//Global 全局应用程序配置
var Global = global.Def

//Services 服务中心
var Services services.IService = services.Def

//Conf 配置组件
var Conf creator.IConf = creator.Conf

//IContext 请求上下文
type IContext = context.IContext

//Component 基础组件
var Component = components.Def
