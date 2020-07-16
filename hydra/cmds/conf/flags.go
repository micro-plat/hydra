package conf

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var coverIfExists = false

//getInstallFlags 获取运行时的参数
func getInstallFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "cover,v",
		Destination: &coverIfExists,
		Usage:       `-覆盖配置，覆盖配置中心和本地服务`,
	})
	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}

//getShowFlags 获取运行时的参数
func getShowFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}
