package main

import (
	"github.com/micro-plat/hydra/examples/rpcserver/services/order"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20-test"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("rpc-api"),
		hydra.WithDebug())

	app.API("/order/query", order.NewQueryHandler)
	app.RPC("/order/bind", order.NewBindHandler)
	app.Start()
}
