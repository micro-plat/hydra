package status

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/global"
	hydracmds "github.com/micro-plat/hydra/hydra/cmds"

	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

var vname string

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "status",
			Usage:  "查询服务状态",
			Flags:  pkgs.GetAppNameFlags(&vname),
			Action: doStatus,
		}
	})
}

func doStatus(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()
	//3.创建本地服务
	hydraSrv, err := hydracmds.GetService(c)
	if err != nil {
		return err
	}
	status, err := hydraSrv.Status()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, status)
}
