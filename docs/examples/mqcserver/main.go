package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra-20"),
		hydra.WithSystemName("collector"),
		hydra.WithServerTypes("mqc-api"),
		hydra.WithDebug())

	app.Conf.MQC.SetSubConf("server", `
	{
		"proto":"redis",
		"addrs":[
				"192.168.0.111:6379",
				"192.168.0.112:6379"
		],
		"db":1,
		"dial_timeout":10,
		"read_timeout":10,
		"write_timeout":10,
		"pool_size":10
}
`)
	app.Conf.MQC.SetSubConf("queue", `{
		"queues":[
			{
				"queue":"hydra:100:0",
				"service":"/message/handle"
			}
		]
	}`)
	app.Conf.Plat.SetVarConf("queue", "queue", `
{
	"proto":"redis",
	"addrs":[
			"192.168.0.111:6379",
			"192.168.0.112:6379"
	],
	"db":1,
	"dial_timeout":10,
	"read_timeout":10,
	"write_timeout":10,
	"pool_size":10
}
`)
	app.Flow("/message/handle", msgHandle)
	app.Micro("/message/send", send)
	app.Start()
}
func msgHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("---------收到消息---------")
	ctx.Log.Info(ctx.Request.GetBody())
	return "success"
}
func send(ctx *context.Context) (r interface{}) {

	queue, err := ctx.GetContainer().GetQueue()
	if err != nil {
		return err
	}
	if err = queue.Push("hydra:100:0", `{"id":"1001"}`); err != nil {
		return err
	}
	return "success"
}
