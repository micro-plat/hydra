package conf

import (
	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:  "conf",
			Usage: "配置管理, 查看、安装配置信息",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "-查看配置，从注册中获取配置并显示",
					Action: showNow,
					Flags:  getShowFlags(),
				},
				{
					Name:   "install",
					Usage:  "-安装配置，将配置信息安装到注册中心",
					Flags:  getInstallFlags(),
					Action: installNow,
				},
			},
		}
	})
}
func showNow(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	//2. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(global.Current().GetRegistryAddr()) == registry.LocalMemory {
		if err := pkgs.Pub2Registry(true); err != nil {
			return err
		}
	}

	//3. 显示配置
	return showConf(global.Current().GetRegistryAddr(),
		global.Current().GetPlatName(),
		global.Current().GetSysName(),
		global.Current().GetServerTypes(),
		global.Current().GetClusterName())
}

func installNow(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	//2.检查是否安装注册中心配置
	if registry.GetProto(global.Current().GetRegistryAddr()) != registry.LocalMemory {
		if err := pkgs.Pub2Registry(coverIfExists); err != nil {
			logs.Log.Error("安装到配置中心:", compatible.FAILED)
			return err
		}
		logs.Log.Info("安装到配置中心:", compatible.SUCCESS)
		return
	}
	return nil

}
