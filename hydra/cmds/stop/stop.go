package stop

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/daemon"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

var vname string

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "stop",
			Usage:  "停止服务",
			Flags:  pkgs.GetAppNameFlags(&vname),
			Action: doStop,
		}
	})
}

func doStop(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()
	service, err := daemon.New(pkgs.GetAppNameDesc(vname))
	if err != nil {
		return err
	}
	msg, err := service.Stop()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
