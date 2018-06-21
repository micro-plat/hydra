package main

import (
	"github.com/micro-plat/hydra/examples/apiserver/services/order"
	"github.com/micro-plat/hydra/examples/apiserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-21"),
		hydra.WithSystemName("collector"),
		//hydra.WithServerTypes("api-web-rpc"),
		hydra.WithServerTypes("web"),
		hydra.WithDebug())
	app.Micro("/user/login", user.NewLoginHandler)
	app.Micro("/order/query", order.NewQueryHandler)
	app.Micro("/order/bind", order.NewBindHandler)

	app.Start()
}
