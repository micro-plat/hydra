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
	app.Cli.Append(hydra.ModeRun, cli.StringFlag{
		Name:  "ip,i",
		Usage: "IP地址",
	})
	app.Cli.Validate(hydra.ModeRun, func(c *cli.Context) error {
		if !c.IsSet("ip") {
			return fmt.Errorf("未设置ip地址")
		}
		return nil
	})
	app.Initializing(func(component.IContainer) error {
		fmt.Println("ip:", app.Cli.Context().String("ip"))
		return nil
	})
	app.Start()
}
func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
