package hydra

import (
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/builder"
	"github.com/micro-plat/hydra/services"
)

//Application 全局应用程序
var Application = application.DefApp

//Services 服务中心
var Services = services.Registry

//Conf 配置组件
var Conf = builder.Conf

//IContext 请求上下文
type IContext = context.IContext

//Component 基础组件
var Component = components.Def
