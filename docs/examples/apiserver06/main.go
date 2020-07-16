package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithClusterName("prod"),
		hydra.WithSystemName("apiserver"),
		hydra.WithConfFlag("port", "api服务器启动端口"),
	)

	app.Cli.Conf.OnStarting(callback)
	app.API("/api", api)
	app.Start()
}
func api(ctx hydra.IContext) interface{} {
	ctx.Log().Info("--------api--------------")
	return map[string]interface{}{
		"name": "colin",
		"id":   ctx.Request().GetInt("id"),
	}
}
func callback(c hydra.ICli) error {
	if !c.IsSet("port") {
		return fmt.Errorf("未设置端口号")
	}
	hydra.Conf.API(c.String("port"))
	return nil
}
