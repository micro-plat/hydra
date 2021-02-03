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
		global.Def.RegistryAddr,
		cover); err != nil {
		return err
	}
	return nil
}

//GetFixedFlags 获取服务名称flags
func GetFixedFlags(isFixed *bool) []cli.Flag {
	flags := make([]cli.Flag, 0, 1)
	flags = append(flags, cli.BoolFlag{
		Name:        "fixed,fd",
		Destination: isFixed,
		Usage:       `-指定服务名称与程序名称相同`,
	})
	return flags

}

//GetBaseFlags 获取运行时的参数
func GetBaseFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 4)
	flags = append(flags, registryFlag)
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
	Usage:       `-注册中心地址。格式：proto://host。如：zk://ip1,ip2  或 fs://../`,
}

var platFlag = cli.StringFlag{
	Name:        "plat,p",
	Destination: &global.FlagVal.PlatName,
	Usage:       "-平台名称",
}

var sysNameFlag = cli.StringFlag{
	Name:        "system,s",
	Destination: &global.FlagVal.SysName,
	Usage:       "-系统名称,默认为当前应用程序名称",
}
var serverTypesFlag = cli.StringFlag{
	Name:        "server-types,S",
	Destination: &global.FlagVal.ServerTypeNames,
	Usage:       fmt.Sprintf("-服务类型，有api,web,rpc,cron,mqc,ws。多个以“-”分割"),
}
var clusterFlag = cli.StringFlag{
	Name:        "cluster,c",
	Destination: &global.FlagVal.ClusterName,
	Usage:       "-集群名称，默认值为：prod",
}
