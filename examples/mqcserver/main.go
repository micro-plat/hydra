package main

import (
	"github.com/micro-plat/hydra/examples/mqcserver/services/order"
	"github.com/micro-plat/hydra/examples/mqcserver/services/user"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("mqc-api-rpc"),
		hydra.WithDebug())
	app.Conf.MQC.SetSubConf("metric", `{
		"host":"http://192.168.0.185:8086",
	"dataBase":"mqcserver",
	"cron":"@every 10s",
	"userName":"",
	"password":""
	}	`)
	app.Conf.API.SetSubConf("metric", `{
		"host":"http://192.168.0.185:8086",
	"dataBase":"mqcserver",
	"cron":"@every 10s",
	"userName":"",
	"password":""
	}	`)
	app.Conf.RPC.SetSubConf("metric", `{
		"host":"http://192.168.0.185:8086",
	"dataBase":"mqcserver",
	"cron":"@every 10s",
	"userName":"",
	"password":""
	}	`)
	app.Flow("/order/query", order.NewQueryHandler)
	app.Flow("/order/bind", order.NewBindHandler)
	app.Micro("/message/send", user.NewLoginHandler)
	app.Micro("/order/bind", order.NewBindHandler)
	app.Start()
}
