package pub

import (
	"github.com/urfave/cli"
)

var runInstall = ""
var pwd string

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 1)
	flags = append(flags, cli.StringFlag{
		Name:        "install,i",
		Destination: &runInstall,
		Usage:       `-按需执行install命令`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "pwd",
		Destination: &pwd,
		Usage:       `-远程服务器密码`,
		Required:    true,
	})
	return flags
}
