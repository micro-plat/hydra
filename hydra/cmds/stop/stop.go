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
			Usage:  "停止服务",
			Action: doStop,
		})
}

func doStop(c *cli.Context) (err error) {

	//关闭日志显示
	application.Current().Log().Pause()
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
