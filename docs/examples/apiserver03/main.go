package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
	)

	app.API("/order/request", hello)
	app.API("/member/login", login)
	app.API("/index", index)
	app.Start()
}
func hello(ctx hydra.IContext) interface{} {
	return "success"
}
func login(ctx hydra.IContext) interface{} {
	ctx.User().Auth().Response(map[string]interface{}{
		"uid": "abc",
	})
	return "success"
}
func index(ctx hydra.IContext) interface{} {
	return ctx.User().Auth().Request()
}
