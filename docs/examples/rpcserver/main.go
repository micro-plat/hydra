package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("qxgrs"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())

	app.API("/test", handle)
	app.Start()
}
func handle(ctx *context.Context) (r interface{}) {
	service := "upchannel/notify/crawl@micro-services.qxgrs"
	ctx.Log.Info("请求服务:", service)
	ctx.Log.Info(ctx.RPC.Request(service, nil, nil, true))
	return "OK"
}
