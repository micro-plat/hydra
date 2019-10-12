package main

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/docs/examples/apiserver/services/order"
	"github.com/micro-plat/hydra/hydra"
)

func main() {

	app := hydra.NewApp(
		hydra.WithPlatName("hydra-780"),
		hydra.WithSystemName("apiserver"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())

	// app.Conf.API.SetMain(conf.NewAPIServerConf(":8098").WithDNS("api.hydra.com"))
	app.Conf.API.SetHeaders(conf.NewHeader().WithCrossDomain())
	app.Micro("/order", order.NewQueryHandler, "/order/tags")
	app.Micro("/hello/Get.do", helloWorld, "/a/bc")

	app.Start()
}
func helloWorld(ctx *context.Context) (r interface{}) {
	ctx.Log.Info(ctx.GetContainer().GetTags(ctx.Service))
	// ctx.Response.SetXML()
	return "hello"
}
