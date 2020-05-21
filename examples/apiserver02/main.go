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

	app.API("/order/request", request)
	app.API("/order", &OrderService{})

	app.OnServerStarting(func(server.IServerConf) error {
		hydra.Application.Log().Info("on.server.starting")
		return nil
	})

	app.OnServerClosing(func(server.IServerConf) error {
		hydra.Application.Log().Info("on.server.closing")
		return nil
	})

	app.OnHandleExecuting(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("on.handle.executing")
		return nil
	})
	app.OnHandleExecuted(func(ctx hydra.IContext) interface{} {
		ctx.Log().Info("on.handle.executed")
		return nil
	})
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	return "request"
}
