package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithSystemName("apiserver"),
	)
	hydra.Conf.API(":8080")
	app.API("/api", api)
	app.Start()
}
func api(ctx hydra.IContext) interface{} {
	return "success"
}
