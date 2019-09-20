package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

type apiserver struct {
	*hydra.MicroApp
}

func main() {
	app := &apiserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("apiserver"),
			hydra.WithServerTypes("api")),
	}

	app.API("/hello", hello)
	app.Conf.API.SetSubConf("metric", `{
		"host":"http://192.168.106.219:8086",
	"dataBase":"mall_apiserver",
	"cron":"@every 10s",
	"userName":"",
	"password":""
    }	`)

	app.Start()

}

func hello(ctx *context.Context) (r interface{}) {
	return "hello world"
}
