package run

import (
	"os"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "run",
			Usage:  "运行服务",
			Flags:  getFlags(),
			Action: doRun,
		}
	})
}

//doRun 服务启动
func doRun(c *cli.Context) (err error) {
	//1.创建本地服务
	hydraSrv, err := pkgs.GetService(c, os.Args[2:]...)
	if err != nil {
		return err
	}
	err = hydraSrv.Run()
	return pkgs.GetCmdsResult(hydraSrv.ServiceName, "Run", err)
}
