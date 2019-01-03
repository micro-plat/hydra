package main

import (
	"fmt"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra"
	"github.com/urfave/cli"
)

func main() {
	app := hydra.NewApp()
	app.Micro("/hello", helloWorld)
	//添加命令行参数
	app.ArgCtx.Append(cli.StringFlag{
		Name:  "ip,i",
		Usage: "IP地址",
	})

	//参数验证
	app.ArgCtx.Validate = func() error {
		if !app.ArgCtx.IsSet("ip") {
			return fmt.Errorf("未指定ip地址")
		}
		return nil
	}

	app.Initializing(func(c component.IContainer) error {
		//获取参数值
		fmt.Println("ip.address:", app.ArgCtx.String("ip"))
		return nil
	})

	app.Start()
}
func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
