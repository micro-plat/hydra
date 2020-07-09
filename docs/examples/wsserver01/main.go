package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务器各种返回结果
func main() {

	app := hydra.NewApp(hydra.WithServerTypes(http.WS))

	//注册服务
	app.WS("/api", api)

	app.Start()
}

func api(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------api--------------")
	return map[string]interface{}{
		"name": "colin",
		"id":   ctx.Request().GetInt("id"),
	}
}
