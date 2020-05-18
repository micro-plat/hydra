package install

import (
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var installRegistry = false
var coverIfExists = false

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "ir,i",
		Destination: &installRegistry,
		Usage:       `-安装注册中心配置`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "cover,c",
		Destination: &coverIfExists,
		Usage:       `-覆盖配置，覆盖配置中心和本地服务`,
	})
	return flags
}
