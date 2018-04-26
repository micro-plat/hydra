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
		hydra.WithAutoCreateConf(),
		hydra.WithDebug())

	app.Autoflow("/user/login", user.NewLoginHandler)
	app.Autoflow("/order/query", order.NewQueryHandler)
	app.Autoflow("/order/bind", order.NewBindHandler)
	app.Micro("/order/bind", order.NewBindHandler)
	app.Start()
}
