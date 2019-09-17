package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),   //平台名称
		hydra.WithSystemName("demo"),   //系统名称
		hydra.WithClusterName("test"),  //集群名称
		hydra.WithServerTypes("api"),   //只启动http api 服务
		hydra.WithRegistry("fs://../"), //使用本地文件系统作为注册中心
		hydra.WithDebug())
	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
