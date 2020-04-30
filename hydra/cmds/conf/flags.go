package conf

import (
	"fmt"

	"github.com/micro-plat/hydra/application"
	"github.com/urfave/cli"
)

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 4)
	if application.RegistryAddr == "" {
		flags = append(flags, cli.StringFlag{
			Name:        "registry,r",
			Destination: &application.RegistryAddr,
			EnvVar:      "hydra_registry",
			Usage: "\033[;31m*\033[0m" + `注册中心地址,必须项。目前支持zookeeper(zk)和本地文件系统(fs)。注册中心用于保存服务启动和运行参数，
	 服务注册与发现等数据，格式:proto://host。proto的取值有zk,fs; host的取值根据不同的注册中心各不同,
	 如zookeeper则为ip地址(加端口号),多个ip用逗号分隔,如:zk://192.168.0.2,192.168.0.107:12181。本地文
	 件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../;  此参数可以通过命令行参数指定，程序指
	 定，也可从环境变量中获取，环境变量名为:`,
		})
	}
	if application.PlatName == "" && application.SysName == "" && len(application.ServerTypes) == 0 && application.ClusterName == "" {
		flags = append(flags, cli.StringFlag{
			Name:        "name,n",
			EnvVar:      "hydra_name",
			Destination: &application.Name,
			Usage: "\033[;31m*\033[0m" + `服务全名，指服务在注册中心的完整名称，该名称是以/分隔的多级目录结构，完整的表示该服务所在平台，系统，服务
	 类型，集群名称，格式：/平台名称/系统名称/服务器类型/集群名称; 平台名称，系统名称，集群名称可以是任意字母
	 下划线或数字，服务器类型则为目前支持的几种服务器类型有:api,web,rpc,mqc,cron,ws。该参数可从环境变量中获取，
	 环境变量名为: `,
		})
	} else {
		if application.PlatName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "plat,p",
				Destination: &application.PlatName,
				EnvVar:      "platName",
				Usage:       "\033[;31m*\033[0m平台名称",
			})
		}
		if application.SysName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "system,s",
				Destination: &application.SysName,
				EnvVar:      "sysName",
				Usage:       "\033[;31m*\033[0m系统名称",
			})
		}
		if len(application.ServerTypes) == 0 {
			flags = append(flags, cli.StringFlag{
				Name:        "server-types,S",
				Destination: &application.ServerTypeNames,
				EnvVar:      "serverType",
				Usage:       fmt.Sprintf("\033[;31m*\033[0m服务类型，目前支持的服务器类型有api,web,rpc,cron,mqc,ws"),
			})
		}
		if application.ClusterName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "cluster,c",
				Destination: &application.ClusterName,
				EnvVar:      "clusterName",
				Usage:       "\033[;31m*\033集群名称",
			})
		}
	}
	return flags
}
