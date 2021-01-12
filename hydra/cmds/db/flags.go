// +build dev

package db

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var dbName = "db"
var skip bool

//getInstallFlags 获取运行时的参数
func getInstallFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "debug,d",
		Destination: &global.FlagVal.IsDebug,
		Usage:       `-调试模式，打印更详细的系统运行日志`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "db",
		Destination: &dbName,
		Usage:       `-数据库节点名,注册中配置的数据库节点名`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "skip",
		Destination: &skip,
		Usage:       `-跳过执行失败的SQL语句`,
	})
	flags = append(flags, global.DBCli.GetFlags()...)
	return flags
}
