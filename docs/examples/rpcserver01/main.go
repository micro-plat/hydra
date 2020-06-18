package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

//服务注册与系统勾子函数
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(rpc.RPC, http.API),
	)
	app.API("/request", request)
	app.RPC("/rpc", rpcRequest)
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	response, err := hydra.Component.RPC().GetRegularRPC().Request(ctx.Context(), "/rpc", map[string]interface{}{
		"id": 2000,
	})
	if err != nil {
		return err
	}
	return response.Result
}

func rpcRequest(ctx hydra.IContext) (r interface{}) {
	return "request" + ctx.Request().GetString("id")
}
