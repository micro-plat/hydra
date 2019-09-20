package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("file"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())
	app.Conf.API.SetMainConf(`{"address":":8080"}`)

	app.Conf.API.SetSubConf("static", `{
		"exts":[".txt"]
	}`)
	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
