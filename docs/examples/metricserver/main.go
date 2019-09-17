package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat_test_9_2"),
		hydra.WithSystemName("demo"),
		hydra.WithServerTypes("api"),
	)

	app.Conf.API.SetMainConf(`{"address":":9067"}`)
	app.Conf.API.SetSubConf("metric", `{
		"host":"http://192.168.106.219:8086",
	"dataBase":"hydra_metrics",
	"cron":"@every 10s",
	"userName":"",
	"password":""
	}	`)

	app.API("/hello", helloWorld)
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	//panic("hello.err")
	return "hello"
}
