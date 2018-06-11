package main

import (
	"github.com/micro-plat/hydra/examples/wsserver/services/order"
	"github.com/micro-plat/hydra/examples/wsserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("ws"),
		hydra.WithDebug())
	app.WS("/user/login", user.NewLoginHandler)
	app.WS("/order/query", order.NewQueryHandler)
	app.WS("/order/bind", order.NewBindHandler)

	app.Start()
}
