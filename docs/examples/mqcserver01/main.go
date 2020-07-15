package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(mqc.MQC, http.API),
		hydra.WithDebug(),
	)
	app.API("/send", send)
	app.MQC("/order/request", show, "order:request")
	app.Start()

}
func send(ctx hydra.IContext) interface{} {
	hydra.MQC.Add("order:request", "/order/request")
	hydra.C.Queue().GetRegularQueue().Push("order:request", `{"id":500}`)
	hydra.MQC.Remove("order:request", "/order/request")
	return "success"
}
func show(ctx hydra.IContext) interface{} {
	ctx.Log().Info("show....")
	ctx.Log().Info("id:", ctx.Request().GetString("id"))
	return "success"
}
