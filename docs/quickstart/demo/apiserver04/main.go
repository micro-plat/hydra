package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp()
	app.API("/hello", hello)
	app.Conf.API.SetMainConf(`{"address":":8091"}`)
	app.Start()
}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
