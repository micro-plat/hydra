package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {

	app := hydra.NewApp(
		hydra.WithServerTypes(http.API, cron.CRON),
	)

	app.API("/hello", hello)

	app.CRON("/hello", hello, "@every 5s")

	app.Start()
}

func hello(ctx hydra.IContext) interface{} {
	return "success"
}
