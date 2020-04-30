package stop

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/hydra/cmds/daemon"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "stop",
			Usage:  "停止服务。通过start启动的服务，使用此命令可停止服务",
			Action: doStop,
		})
}

func doStop(c *cli.Context) (err error) {
	service, err := daemon.New(application.AppName, application.AppName)
	if err != nil {
		return err
	}
	msg, err := service.Stop()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
