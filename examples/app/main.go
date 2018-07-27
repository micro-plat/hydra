package main

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp()
	app.Micro("/hello", helloWorld)
	app.Start()
}
func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
