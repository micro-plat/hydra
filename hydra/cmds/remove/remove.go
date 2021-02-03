package remove

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"

	"github.com/urfave/cli"
)

var isFixed bool

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "remove",
			Usage:  "删除服务，从本地服务器移除服务",
			Flags:  pkgs.GetFixedFlags(&isFixed),
			Action: doRemove,
		}
	})
}
func doRemove(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()

	//3.创建本地服务
	hydraSrv, err := pkgs.GetService(c, isFixed)
	if err != nil {
		return err
	}
	err = hydraSrv.Uninstall()
	return pkgs.GetCmdsResult(hydraSrv.DisplayName, "Remove", err)

}
