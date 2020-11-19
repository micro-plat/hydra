package start

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"

	"github.com/micro-plat/hydra/global"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "start",
			Usage:  "启动服务",
			Action: doStart,
		}
	})
}

func doStart(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()

	//3.创建本地服务
	hydraSrv, err := pkgs.GetService(c)
	if err != nil {
		return err
	}
	err = hydraSrv.Start()
	return pkgs.GetCmdsResult(hydraSrv.ServiceName, "Start", err)
}
