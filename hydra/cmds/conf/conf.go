package conf

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
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
	if installRegistry {
		if err := pkgs.Pub2Registry(coverIfExists); err != nil {
			logs.Log.Errorf("安装到配置中心:", pkgs.Failed)
			return err
		}
		logs.Log.Debug("安装到配置中心:" + pkgs.Success)
	}

	//3. 显示配置
	return show()
}
