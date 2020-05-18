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
			Usage:  "启动服务",
			Action: doStart,
		})
}

func doStart(c *cli.Context) (err error) {

	//关闭日志显示
	application.Current().Log().Pause()
	service, err := daemon.New(application.DefApp.GetLongAppName(), application.Usage)
	if err != nil {
		return err
	}
	msg, err := service.Start()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
