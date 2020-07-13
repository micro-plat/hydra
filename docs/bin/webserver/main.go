package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.Web),
		hydra.WithDebug(),
	)
	hydra.Conf.Web(":81")
	app.Micro("/order/request", hello)
	app.Start()
}

func hello(ctx hydra.IContext) interface{} {
	return "success"
}
