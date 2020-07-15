package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithClusterName("prod"),
	)
	app.API("/api", api)
	app.Start()
}
func api(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------api--------------")
	return map[string]interface{}{
		"name": "colin",
		"id":   ctx.Request().GetInt("id"),
	}
}
