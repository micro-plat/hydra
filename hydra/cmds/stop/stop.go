package stop

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
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
	//3.创建本地服务
	hydraSrv, err := pkgs.GetService(c)
	if err != nil {
		return err
	}

	err = hydraSrv.Stop()
	return pkgs.GetCmdsResult(hydraSrv.ServiceName, "Stop", err)
}
