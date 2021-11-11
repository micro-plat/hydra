package conf

import (
	"github.com/lib4dev/cli/cmds"
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
				{
					Name:   "encrypt",
					Usage:  "-使用内置加密方法加密配置数据",
					Action: encrypt,
					Flags:  getEncryptFlags(),
				},
				{
					Name:   "export",
					Usage:  "-使用内置加密方法加密配置数据",
					Action: exportNow,
					Flags:  getExportFlags(),
				},
			},
		}
	})
}
