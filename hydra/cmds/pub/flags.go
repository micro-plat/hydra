package pub

import (
	"github.com/urfave/cli"
)

var runRun = false
var runStart = false
var runInstall = ""
var pwd string

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 1)
	flags = append(flags, cli.BoolFlag{
		Name:        "start,s",
		Destination: &runStart,
		Usage:       `-按需执行start命令`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "install,i",
		Destination: &runInstall,
		Usage:       `-按需执行install命令`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "run",
		Destination: &runRun,
		Usage:       `-按需执行run命令`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "pwd",
		Destination: &pwd,
		Usage:       `-远程服务器密码`,
		Required:    true,
	})
	return flags
}
