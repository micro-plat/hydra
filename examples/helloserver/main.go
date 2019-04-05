package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("helloserver"),

		hydra.WithDebug(),
	)

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
