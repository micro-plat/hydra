package main

import (
	"github.com/micro-plat/hydra/examples/apiserver/services/order"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra"),
		hydra.WithSystemName("apiserver"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())
	app.Micro("/order/query", order.NewQueryHandler)
	app.Micro("/order", order.NewBindHandler)
	app.Start()
}
