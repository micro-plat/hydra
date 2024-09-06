package install

import (
	"os"

	"github.com/lib4dev/cli/cmds"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var isFixed bool

func init() {
	cmds.RegisterFunc(func() cli.Command {
		flags := pkgs.GetFixedFlags(&isFixed)
		flags = append(flags, getFlags()...)
		return cli.Command{
			Name:   "install",
			Usage:  "安装服务，以服务方式安装到本地系统",
			Flags:  flags,
			Action: doInstall,
		}
	})
}

func doInstall(c *cli.Context) (err error) {

	//1.检查是否有管理员权限
	global.Current().Log().Pause()
	if err = compatible.CheckPrivileges(); err != nil {
		return err
	}

	//2. 绑定应用程序参数
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}
	args := []string{"run", "--nostd"}
	args = append(args, os.Args[2:]...)

	//执行本地处理服务
	handlers := global.Installer.Local.GetHandlers()
	for _, v := range handlers {
		if err := v(); err != nil {
			return err
		}
	}

	//3.创建本地服务
	hydraSrv, err := pkgs.GetService(c, isFixed, args...)
	if err != nil {
		return err
	}
	if coverIfExists {
		hydraSrv.Uninstall()
	}

	err = hydraSrv.Install()
	return pkgs.GetCmdsResult(hydraSrv.DisplayName, "Install", err)
}
