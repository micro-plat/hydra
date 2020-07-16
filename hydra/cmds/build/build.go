package build

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:  "build",
			Usage: "构建应用配置",
			Subcommands: []cli.Command{
				{
					Name:   "updater",
					Usage:  "构建updater文件包",
					Action: createUpdater,
					Flags:  getFlags(),
				},
			},
		}
	})
}
