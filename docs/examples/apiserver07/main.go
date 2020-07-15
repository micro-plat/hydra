package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		// hydra.WithClusterName("prod"),
	)
	hydra.Conf.API(":8888").Gray(`{{$id := get_req "id"}}{{if eq $id "100"}}true{{end}}`, "prod")
	app.API("/api", api)
	app.Start()
}
func api(ctx hydra.IContext) interface{} {
	return "success"
}
