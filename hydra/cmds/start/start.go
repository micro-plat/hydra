package start

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
			Name:   "start",
			Usage:  "启动服务。后台运行服务，日志存入本地文件或日志中心。异常退出或服务器重启会自动启动",
			Action: doStart,
		})
}

func doStart(c *cli.Context) (err error) {
	service, err := daemon.New(application.AppName, application.AppName)
	if err != nil {
		return err
	}
	msg, err := service.Start()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
