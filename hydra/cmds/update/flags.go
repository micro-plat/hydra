package update

import (
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var url string
var coverIfExists = false

//getInstallFlags 获取运行时的参数
func getInstallFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "url,u",
		Required:    true,
		Destination: &url,
		Usage:       `应用下载地址`,
	})
	return flags
}

//getBuildFlags 获取运行时的参数
func getBuildFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 1) // pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "url,u",
		Required:    true,
		Destination: &url,
		Usage:       `应用下载地址`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "cover,v",
		Destination: &coverIfExists,
		Usage:       `-文件已存在是否删除`,
	})
	return flags
}
