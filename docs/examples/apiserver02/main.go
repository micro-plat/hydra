package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server"
)

//服务注册与系统勾子函数
func main() {
	app := hydra.NewApp()

	app.API("/order/request", request, "/order/*")
	app.API("/order", &OrderService{})

	app.OnStarting(func(server.IServerConf) error {
		hydra.Global.Log().Info("server.OnServerStarting")
		return nil
	})

	app.OnClosing(func(server.IServerConf) error {
		hydra.Global.Log().Info("server.OnServerClosing")
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
