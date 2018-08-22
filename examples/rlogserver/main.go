package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("myplat_test_rlogger"),
		hydra.WithSystemName("demo"),
		hydra.WithServerTypes("api"),
		hydra.WithRemoteLogger(),
		hydra.WithDebug())

	app.Conf.API.SetMainConf(`{"address":":7890"}`)
	app.Conf.Plat.SetVarConf("global", "logger", `{
    "level":"All",
		"interval":"1s",
		"layout":{"name":"%name","time":"%datetime","content":"%content","level":"%l","session":"%session"},
    "service":"/hydra/log/save@logging.logging_debug"
}`)
	app.Conf.API.SetSubConf("metric", `{
	"host":"http://192.168.0.185:8086",
	"dataBase":"rlogserver",
	"cron":"@every 10s",
	"userName":"",
	"password":""
	}`)

	app.Micro("/hello", helloWorld)
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
