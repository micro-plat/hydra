package main

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/hydra"
)

var ch chan *conf.Queue

func main() {
	ch = make(chan *conf.Queue, 2)
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
	app.SetMQCDynamicQueue(ch)
	app.Flow("/message/handle", msgHandle)
	app.Micro("/message/send", send)
	app.Micro("/consume", consume)
	app.Micro("/unconsume", unconsume)
	app.Start()
}
func msgHandle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("---------收到消息---------")
	ctx.Log.Info(ctx.Request.GetBody())
	return "success"
}
func consume(ctx *context.Context) (r interface{}) {
	if err := ctx.Request.Check("name"); err != nil {
		return err
	}
	q := &conf.Queue{Queue: ctx.Request.GetString("name"), Service: "/message/handle"}
	ch <- q
	return q
}
func unconsume(ctx *context.Context) (r interface{}) {
	if err := ctx.Request.Check("name"); err != nil {
		return err
	}
	q := &conf.Queue{Queue: ctx.Request.GetString("name"), Disable: true}
	ch <- q
	return q
}

func send(ctx *context.Context) (r interface{}) {
	if err := ctx.Request.Check("name"); err != nil {
		return err
	}

	queue, err := ctx.GetContainer().GetQueue()
	if err != nil {
		return err
	}
	if err = queue.Push(ctx.Request.GetString("name"), `{"id":"1001"}`); err != nil {
		return err
	}
	return "success"
}
