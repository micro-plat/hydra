package conf

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "conf",
			Usage:  "查看配置信息",
			Flags:  getFlags(),
			Action: doConf,
		}
	})
}

func doConf(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}

	//2.检查是否安装注册中心配置
	if installRegistry && registry.GetProto(global.Current().GetRegistryAddr()) != registry.LocalMemory {
		if err := pkgs.Pub2Registry(coverIfExists); err != nil {
			logs.Log.Error("安装到配置中心:", pkgs.Failed)
			return err
		}
		logs.Log.Info("安装到配置中心:" + pkgs.Success)
		return
	}

	//3. 处理本地内存作为注册中心的服务发布问题
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
