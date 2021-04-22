package conf

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

var coverIfExists = false
var importConf string

//getInstallFlags 获取运行时的参数
func getInstallFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "cover,v",
		Destination: &coverIfExists,
		Usage:       `-覆盖配置，覆盖配置中心和本地服务`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "debug,d",
		Destination: &global.FlagVal.IsDebug,
		Usage:       `-调试模式，打印更详细的系统运行日志，避免将详细的错误信息返回给调用方`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "import",
		Destination: &importConf,
		Usage:       `-导入配置文件`,
	})
	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}

var extNode string

//getShowFlags 获取运行时的参数
func getShowFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "node",
		Destination: &extNode,
		Usage:       `-扩展节点名称`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "debug,d",
		Destination: &global.FlagVal.IsDebug,
		Usage:       `-调试模式，打印更详细的系统运行日志，避免将详细的错误信息返回给调用方`,
	})
	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}

var orgData string

//getEncryptFlags 获取运行时的参数
func getEncryptFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "data",
		Destination: &orgData,
		Usage:       `-需要加密数据`,
		Required:    true,
	})

	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}

var coverConfIfExists = false
var confEncrypt = false
var confExportPath string

//getExportFlags 获取导出配置时的参数
func getExportFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.BoolFlag{
		Name:        "cover,v",
		Destination: &coverConfIfExists,
		Usage:       `-导出配置文件已存在是否删除`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "debug,d",
		Destination: &global.FlagVal.IsDebug,
		Usage:       `-调试模式，打印更详细的系统运行日志，避免将详细的错误信息返回给调用方`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "encrypt,e",
		Destination: &confEncrypt,
		Usage:       `-导出配置是否进行加密`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "out,o",
		Destination: &confExportPath,
		Usage:       `-配置文件导出地址`,
	})
	flags = append(flags, global.ConfCli.GetFlags()...)
	return flags
}
