package main

import (
	"time"

	"github.com/micro-plat/hydra"

	"github.com/micro-plat/hydra/components/pkgs/apm/apmtypes"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//服务器各种返回结果
func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
		hydra.WithAPM(apmtypes.SkyWalking),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("apiserver03"),
	)

	app.API("/order/request", hello, router.WithEncoding("gbk"))
	app.API("/member/login", login)
	app.API("/index", index)
	app.Start()
}
func hello(ctx hydra.IContext) interface{} {
	return "order/request"
}
func login(ctx hydra.IContext) interface{} {
	ctx.User().Auth().Response(map[string]interface{}{
		"uid": "abc",
	})
	return "success"
}
func index(ctx hydra.IContext) interface{} {
	time.Sleep(time.Second * 4)
	return hydra.C.UUID().ToString()
}
