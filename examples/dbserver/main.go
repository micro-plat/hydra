package main

import (
	"github.com/micro-plat/hydra/examples/dbserver/services/order"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydrav-db"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("api"),
		hydra.WithAutoCreateConf(),
		hydra.WithDebug())
	app.Micro("/order/query", order.NewQueryHandler)
	app.Micro("/order/request", order.NewRequestHandler)
	app.Start()
}
