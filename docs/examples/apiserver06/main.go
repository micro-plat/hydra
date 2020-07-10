package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
	app.API("/api", api)
	hydra.Conf.API(":8080").Limit(limiter.NewRule("/**", 1, limiter.WithReponse(302, "http://www.baidu.com")))

	app.Start()
}
func api(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------api--------------")
	return map[string]interface{}{
		"name": "colin",
		"id":   ctx.Request().GetInt("id"),
	}
}
