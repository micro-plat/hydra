package remove

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
			Name:   "remove",
			Usage:  "删除服务",
			Flags:  pkgs.GetAppNameFlags(&vname),
			Action: doRemove,
		}
	})
}
func doRemove(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()

	//3.创建本地服务
	hydraSrv, err := hydracmds.GetService(c)
	if err != nil {
		return err
	}
	err = hydraSrv.Uninstall()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, err)
}
