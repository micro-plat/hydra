package conf

import (
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/registry"
	"github.com/urfave/cli"
)

func installNow(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return err
	}

	//2.检查是否安装注册中心配置
	if registry.GetProto(global.Current().GetRegistryAddr()) != registry.LocalMemory {

		//导入配置
		input, err1 := creator.GetImportConfs(importConf)
		if err1 != nil {
			logs.Log.Error("导入配置到配置中心:", compatible.FAILED)
			return err1
		}

		if err := pkgs.Pub2Registry(coverIfExists, input); err != nil {
			logs.Log.Error("安装到配置中心:", compatible.FAILED)
			return err
		}
		logs.Log.Info("安装到配置中心:", compatible.SUCCESS)
		return
	}

	return nil
}
