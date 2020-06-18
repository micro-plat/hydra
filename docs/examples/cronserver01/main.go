package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(cron.CRON),
		hydra.WithDebug(),
	)
	app.CRON("/order/request", request, "@every 5s")
	app.CRON("/order/query", query)
	app.Start()
}
func request(ctx hydra.IContext) interface{} {
	hydra.CRON.Add("@now", "/order/query")
	return "success"
}
func query(ctx hydra.IContext) interface{} {
	return "success"
}
