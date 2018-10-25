package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat_test_9_1"),
		hydra.WithSystemName("demo"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug(),
		//	hydra.WithRegistry("zk://192.168.0.107"),
		//	hydra.WithClusterName("test")
	)

	app.Conf.API.SetMainConf(`{"address":":9067"}`)

	app.Micro("/hello", helloWorld)
	app.Start()
}

type Input struct {
	Id int `json:"id" form:"id"`
}

func helloWorld(ctx *context.Context) (r interface{}) {
	var input Input
	ctx.Request.Bind(&input)
	return &input
}
