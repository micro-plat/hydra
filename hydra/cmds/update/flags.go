package update

import (
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var url string

//getInstallFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "url,u",
		Required:    true,
		Destination: &url,
		Usage:       "\033[;31m*\033[0m" + `应用下载地址`,
	})
	return flags
}
