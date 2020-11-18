package start

import (
	"fmt"
	"github.com/micro-plat/cli/cmds"
	hydracmds "github.com/micro-plat/hydra/hydra/cmds"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/errs"
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
	hydraSrv, err := hydracmds.GetService(c)
	if err != nil {
		fmt.Println("doStart.hydracmds.GetService:", err)
		return err
	}
	err = hydraSrv.Start()
	if err != nil {
		fmt.Println("doStart. hydraSrv.Start:", err)
		return err
	}
	return errs.NewIgnoreError(0, err)
}
