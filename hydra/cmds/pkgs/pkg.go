package pkgs

import (
	"fmt"

	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/urfave/cli"
)

//Pub2Registry 发布到注册中心
func Pub2Registry(cover bool) error {

	//2.发布到配置中心
	if err := creator.Conf.Pub(global.Current().GetPlatName(),
		global.Current().GetSysName(),
		global.Current().GetClusterName(),
		global.Def.RegistryAddr, cover); err != nil {
		return err
	}
	return nil
}

//GetAppNameFlags 获取服务名称flags
func GetAppNameFlags(vname *string) []cli.Flag {
	flags := make([]cli.Flag, 0, 1)
	flags = append(flags, cli.StringFlag{
		Name:        "name,n",
		Destination: vname,
		Usage:       `-指定服务名称`,
	})
	return flags

}

//GetAppNameDesc 获取应用程序名称
func GetAppNameDesc(vname string) (string, string) {
	if vname != "" {
		return global.Def.GetLongAppName(vname), global.Def.GetLongAppName(vname)
	}
	return global.Def.GetLongAppName(), global.Usage
}

//GetBaseFlags 获取运行时的参数
func GetBaseFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 4)
	// if global.Def.RegistryAddr == "" {
	// 	flags = append(flags, registryFlag)
	// }
	// if global.Def.Name == "" &&
	// 	global.Def.PlatName == "" && global.Def.SysName == "" && len(global.Def.ServerTypes) == 0 &&
	// 	global.Def.ClusterName == "" {
	// 	flags = append(flags, nameFlag)
	// } else {
	// 	if global.Def.PlatName == "" {
	// 		flags = append(flags, platFlag)
	// 	}
	// 	if global.Def.SysName == "" {
	// 		flags = append(flags, sysNameFlag)
	// 	}
	// 	if len(global.Def.ServerTypes) == 0 {
	// 		flags = append(flags, serverTypesFlag)
	// 	}
	// 	if global.Def.ClusterName == "" {
	// 		flags = append(flags, clusterFlag)
	// 	}
	// }

	flags = append(flags, registryFlag)
	flags = append(flags, nameFlag)
	flags = append(flags, platFlag)
	flags = append(flags, sysNameFlag)
	flags = append(flags, serverTypesFlag)
	flags = append(flags, clusterFlag)
	return flags
}

var registryFlag = cli.StringFlag{
	Name:        "registry,r",
	Destination: &global.FlagVal.RegistryAddr,
	EnvVar:      "registry",
	//Value:       "lm://.",
	Usage: "\033[;31m*\033[0m" + `注册中心地址。目前支持zookeeper(zk)和本地文件系统(fs)。注册中心用于保存服务启动和运行参数，
	服务注册与发现等数据，格式:proto://host。proto的取值有zk,fs; host的取值根据不同的注册中心各不同,
	如zookeeper则为ip地址(加端口号),多个ip用逗号分隔,如:zk://192.168.0.2,192.168.0.107:12181。本地文
	件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../;  此参数可以通过命令行参数指定，程序指
	定，也可从环境变量中获取，环境变量名为:`,
}
var nameFlag = cli.StringFlag{
	Name:        "name,n",
	EnvVar:      "name",
	Destination: &global.FlagVal.Name,
	Usage: "\033[;31m*\033[0m" + `服务全名，指服务在注册中心的完整名称，该名称是以/分隔的多级目录结构，完整的表示该服务所在平台，系统，服务
	类型，集群名称，格式：/平台名称/系统名称/服务器类型/集群名称; 平台名称，系统名称，集群名称可以是任意字母
	下划线或数字，服务器类型则为目前支持的几种服务器类型有:api,web,rpc,mqc,cron,ws。该参数可从环境变量中获取，
	环境变量名为: `,
}
var platFlag = cli.StringFlag{
	Name:        "plat,p",
	Destination: &global.FlagVal.PlatName,
	Usage:       "\033[;31m*\033[0m平台名称",
}

var sysNameFlag = cli.StringFlag{
	Name:        "system,s",
	Destination: &global.FlagVal.SysName,
	Usage:       "\033[;31m*\033[0m系统名称",
}
var serverTypesFlag = cli.StringFlag{
	Name:        "server-types,S",
	Destination: &global.FlagVal.ServerTypeNames,
	Usage:       fmt.Sprintf("\033[;31m*\033[0m服务类型，目前支持的服务器类型有api,web,rpc,cron,mqc,ws"),
}
var clusterFlag = cli.StringFlag{
	Name:        "cluster,c",
	Destination: &global.FlagVal.ClusterName,
	Usage:       "\033[;31m*\033[0m集群名称",
}
