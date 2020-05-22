package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/hydra/servers/http"
)

//服务注册与系统勾子函数
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
	)

	app.API("/order/request", request, "/order/*")
	app.API("/order", &OrderService{})

	app.OnServerStarting(func(server.IServerConf) error {
		hydra.Application.Log().Info("server.OnServerStarting")
		return nil
	})

	app.OnServerClosing(func(server.IServerConf) error {
		hydra.Application.Log().Info("server.OnServerClosing")
		return nil
	})

	app.OnHandleExecuting(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("global.OnHandleExecuting")
		return nil
	})
	app.OnHandleExecuted(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("global.OnHandleExecuted")
		return nil
	})
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	return "request"
}
