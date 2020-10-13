package main

import (
	"time"

	"github.com/micro-plat/hydra"

	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("apiserver04"),
	)

	app.API("/order/request", hello)
	app.Start()
}
func hello(ctx hydra.IContext) interface{} {
	time.Sleep(time.Second * 2)
	return "order/request"
}
