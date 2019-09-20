package main

import (
	"github.com/micro-plat/hydra/examples/webserver/services/order"
	"github.com/micro-plat/hydra/examples/webserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("web"),
		hydra.WithDebug())

	app.Page("/user/login", user.NewLoginHandler)
	app.Page("/order/query", order.NewQueryHandler)
	app.Micro("/order/bind", order.NewBindHandler)

	app.Start()
}
