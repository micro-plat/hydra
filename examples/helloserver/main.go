package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat_test_9_1"),
		hydra.WithSystemName("demo"),
		hydra.WithServerTypes("api"),
		//	hydra.WithRegistry("zk://192.168.0.107"),
		//	hydra.WithClusterName("test"),
		hydra.WithDebug())

	app.Conf.API.SetMainConf(`{"address":"#address"}`)

	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
