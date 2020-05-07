package install

import (
	"os"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/hydra/cmds/daemon"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "install",
			Usage:  "安装本地服务和远程配置服务",
			Flags:  getFlags(),
			Action: doInstall,
		})
}

func doInstall(c *cli.Context) (err error) {

	//1. 绑定应用程序参数
	if err := application.DefApp.Bind(); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	service, err := daemon.New(application.AppName, application.AppName)
	if err != nil {
		return err
	}
	msg, err := service.Install(os.Args[2:]...)
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
