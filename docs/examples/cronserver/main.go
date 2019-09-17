package main

import (
	"github.com/micro-plat/hydra/examples/cronserver/services/order"
	"github.com/micro-plat/hydra/examples/cronserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("cron-rpc"),
		hydra.WithDebug())

	app.Flow("/user/login", user.NewLoginHandler)
	app.Flow("/order/query", order.NewQueryHandler)
	app.Flow("/order/bind", order.NewBindHandler)
	app.Micro("/order/bind", order.NewBindHandler)
	app.Start()
}
