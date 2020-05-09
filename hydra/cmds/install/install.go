package install

import (
	"os"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/hydra/cmds/daemon"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf/builder"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "install",
			Usage:  "安装服务。将配置信息安装到注册中心，并在本地创建服务。安装完成后可通过'start'命令启动服务",
			Flags:  getFlags(),
			Action: doInstall,
		})
}

func doInstall(c *cli.Context) (err error) {

	//1.检查是否有管理员权限
	application.Current().Log().Pause()
	if err = application.CheckPrivileges(); err != nil {
		return err
	}
	//2. 绑定应用程序参数
	if err := application.DefApp.Bind(); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}
	//3.检查是否只安装本地服务
	if !onlyInstallLocalService {

		//1. 加载配置信息
		if err := builder.Conf.Load(); err != nil {
			return err
		}

		//2. 发布到配置中心
		r, err := registry.NewRegistry(application.Current().GetRegistryAddr(), application.Current().Log())
		if err != nil {
			return err
		}
		if err := builder.Conf.Pub(application.Current().GetPlatName(), application.Current().GetSysName(), application.Current().GetClusterName(), r); err != nil {
			return err
		}
	}

	//4.创建本地服务
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
