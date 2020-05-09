package install

import (
	"github.com/micro-plat/hydra/hydra/cmds/conf"
	"github.com/urfave/cli"
)

var onlyInstallLocalService = false

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := conf.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "local,l",
		Destination: &onlyInstallLocalService,
		Usage:       `-只安装本地服务`,
	})
	return flags
}
