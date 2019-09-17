package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat"),
		hydra.WithSystemName("helloserver"),
		hydra.WithServerTypes("api"), //服务器类型为http api
		hydra.WithDebug(),
	)

	app.Micro("/hello", helloWorld)
	app.Start()
}

type result struct {
	Name string `json:"name" xml:"name"`
}

func helloWorld(ctx *context.Context) (r interface{}) {
	// ctx.Response.SetXML()
	return &result{Name: "hello"}
}
