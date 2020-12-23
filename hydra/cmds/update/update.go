package update

import (
	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/global"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:  "update",
			Usage: "更新应用，将服务发布到远程服务器",
			Subcommands: []cli.Command{
				{
					Name:   "install",
					Usage:  "更新当前应用。下载安装包并自动安装",
					Flags:  getInstallFlags(),
					Action: doUpdate,
				}, {
					Name:   "build",
					Usage:  "打包安装包。创建压缩包，生成安装配置",
					Flags:  getBuildFlags(),
					Action: doBuild,
				},
			},
		}
	})
}

func doUpdate(c *cli.Context) (err error) {
	//1. 绑定应用程序参数
	global.Current().Log().Pause()
	if err := global.Def.Bind(c); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}

	//2.获取应用服务器包配置信息
	pkg, err := GetPackage(url)
	if err != nil {
		return err
	}

	//3. 检查是否需要更新
	if ok, err := pkg.Check(); !ok {
		logs.Log.Errorf("更新失败：%v", err)
		return err
	}

	//4.立即更新
	if err = pkg.Update(logs.Log, func() {
		//关闭服务器
	}); err != nil {
		return err
	}
	return nil
}
