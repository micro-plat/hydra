package main

import (
	"github.com/micro-plat/hydra/examples/wsserver/services/order"
	"github.com/micro-plat/hydra/examples/wsserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-22"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("api-ws"),
		hydra.WithDebug())
	app.WS("/auth/login", user.NewLoginHandler)
	app.WS("/order/query", order.NewQueryHandler)
	app.WS("/order/bind", order.NewBindHandler)
	app.Conf.WS.SetSubConf("auth", `
		{
			"jwt": {
				"exclude": ["/auth/login"],
				"expireAt": 36000,
				"mode": "HS512",
				"name": "__jwt__",
				"redirect":"/auth/login",
				"secret": "12345678"
			}
		}
		`)
	app.Start()
}
