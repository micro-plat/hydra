package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),        //平台名称
		hydra.WithSystemName("helloserver"), //系统名称
		hydra.WithDebug())

	app.Micro("/hello", hello)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
