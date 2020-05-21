package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/servers/http"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
	)

	// app.API("/order/request/:tp", request)
	app.Start()
}
