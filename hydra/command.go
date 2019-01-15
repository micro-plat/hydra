package hydra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

//VERSION 版本号
var VERSION = "2.0.0"

func (m *MicroApp) getCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Version = VERSION
	app.Usage = "基于hydra的微服务应用"
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help,h",
		Usage: "查看帮助信息",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,v",
		Usage: "查看版本信息",
	}
	app.Commands = m.getCommands()
	return app
}

func (m *MicroApp) getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "run",
			Usage:  "运行服务",
			Flags:  m.getStartFlags(),
			Action: m.action,
		}, {
			Name:   "start",
			Usage:  "启动服务",
			Action: m.startAction,
		},
		{
			Name:   "stop",
			Usage:  "停止服务",
			Action: m.stopAction,
		}, {
			Name:   "install",
			Usage:  "安装注册中心配置和本地服务",
			Flags:  m.getStartFlags(),
			Action: m.installAction,
		}, {
			Name:   "registry",
			Usage:  "注册中心配置",
			Flags:  m.getStartFlags(),
			Action: m.registryAction,
		}, {
			Name:   "service",
			Usage:  "本地服务配置",
			Flags:  m.getStartFlags(),
			Action: m.serviceAction,
		},
		{
			Name:   "remove",
			Usage:  "删除服务",
			Action: m.removeAction,
		}, {
			Name:   "status",
			Usage:  "查询服务状态",
			Action: m.statusAction,
		}, {
			Name:   "conf",
			Usage:  "查看配置信息",
			Flags:  m.getStartFlags(),
			Action: m.queryConfigAction,
		},
	}
}
func (m *MicroApp) getStartFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 4)
	if m.RegistryAddr == "" {
		flags = append(flags, cli.StringFlag{
			Name:        "registry,r",
			Destination: &m.RegistryAddr,
			EnvVar:      "hydra_registry",
			Usage: "\033[;31m*\033[0m" + `注册中心地址,必须项。目前支持zookeeper(zk)和本地文件系统(fs)。注册中心用于保存服务启动和运行参数，
	 服务注册与发现等数据，格式:proto://host。proto的取值有zk,fs; host的取值根据不同的注册中心各不同,
	 如zookeeper则为ip地址(加端口号),多个ip用逗号分隔,如:zk://192.168.0.2,192.168.0.107:12181。本地文
	 件系统为本地文件路径，可以是相对路径或绝对路径,如:fs://../;  此参数可以通过命令行参数指定，程序指
	 定，也可从环境变量中获取，环境变量名为:`,
		})
	}
	if m.PlatName == "" && m.SystemName == "" && len(m.ServerTypes) == 0 && m.ClusterName == "" {
		flags = append(flags, cli.StringFlag{
			Name:        "name,n",
			EnvVar:      "hydra_name",
			Destination: &m.Name,
			Usage: "\033[;31m*\033[0m" + `服务全名，指服务在注册中心的完整名称，该名称是以/分隔的多级目录结构，完整的表示该服务所在平台，系统，服务
	 类型，集群名称，格式：/平台名称/系统名称/服务器类型/集群名称; 平台名称，系统名称，集群名称可以是任意字母
	 下划线或数字，服务器类型则为目前支持的几种服务器类型有:api,web,rpc,mqc,cron,ws。该参数可从环境变量中获取，
	 环境变量名为: `,
		})
	} else {
		if m.PlatName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "plat,p",
				Destination: &m.PlatName,
				EnvVar:      "hydra_plat",
				Usage:       "\033[;31m*\033[0m平台名称",
			})
		}
		if m.SystemName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "system,s",
				Destination: &m.SystemName,
				EnvVar:      "hydra_system",
				Usage:       "\033[;31m*\033[0m系统名称",
			})
		}
		if len(m.ServerTypes) == 0 {
			flags = append(flags, cli.StringFlag{
				Name:        "server-types,S",
				Destination: &m.ServerTypeNames,
				EnvVar:      "hydra_server_types",
				Usage:       fmt.Sprintf("\033[;31m*\033[0m服务类型，目前支持的服务器类型有%v", supportServerType),
			})
		}
		if m.ClusterName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "cluster,c",
				Destination: &m.ClusterName,
				EnvVar:      "hydra_cluster",
				Usage:       "\033[;31m*\033集群名称",
			})
		}
	}

	flags = append(flags, cli.StringFlag{
		Name:        "trace,t",
		Destination: &m.Trace,
		EnvVar:      "hydra_trace",
		Usage: `-性能跟踪，可选项。用于生成golang的pprof的性能分析数据,支持的模式有:cpu,mem,block,mutex,web。其中web是以http
	 服务的方式提供pprof数据。该参数可从环境变量中获取，环境变量名为:`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "rlog,l",
		Destination: &m.remoteLogger,
		EnvVar:      "hydra_rlog",
		Usage: `-启用远程日志,可选项。默认不启用。指定此参数后启动时自动从注册中心拉取远程日志配置数据，根据配置自动将本地日志以json格
    	 式压缩后定时发送到远程服务器。该参数可从环境变量中获取，环境变量名为:`,
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "rs,R",
		Destination: &m.RemoteQueryService,
		EnvVar:      "hydra-rs",
		Usage: `-启用远程服务,可选项。默然不启用。启用后本地将自动启动一个http服务器。可通过http://host/server/query查询
	 服务状态，通过http://host/update/:version远程更新系统，执行远程更新后服务器将自动从注册中心下载安装包，自动安装并重启服务。
	 该参数可从环境变量中获取，环境变量名为:`,
	})
	flags = append(flags, m.ArgCtx.cmds...)
	return flags
}
