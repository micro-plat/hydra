package hydra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

//VERSION 版本号
var VERSION string = "2.0.0"

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
			Usage:  "立即运行",
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
			Usage:  "安装服务",
			Flags:  m.getStartFlags(),
			Action: m.installAction,
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
			Name:   "config",
			Usage:  "查询配置信息",
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
			Usage:       "注册中心:格式:proto://addr1,addr2",
		})
	}
	if m.PlatName == "" && m.SystemName == "" && len(m.ServerTypes) == 0 && m.ClusterName == "" {
		flags = append(flags, cli.StringFlag{
			Name:        "name,n",
			EnvVar:      "hydra_name",
			Destination: &m.Name,
			Usage:       "服务全称:格式:/平台名称/系统名称/服务器类型/集群名称",
		})
	} else {
		if m.PlatName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "plat,p",
				Destination: &m.PlatName,
				EnvVar:      "hydra_plat",
				Usage:       "平台名称",
			})
		}
		if m.SystemName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "system,s",
				Destination: &m.SystemName,
				EnvVar:      "hydra_system",
				Usage:       "系统名称",
			})
		}
		if len(m.ServerTypes) == 0 {
			flags = append(flags, cli.StringFlag{
				Name:        "server-types,s",
				Destination: &m.ServerTypeNames,
				EnvVar:      "hydra_server_types",
				Usage:       fmt.Sprintf("服务类型%v", supportServerType),
			})
		}
		if m.ClusterName == "" {
			flags = append(flags, cli.StringFlag{
				Name:        "cluster,c",
				Destination: &m.ClusterName,
				EnvVar:      "hydra_cluster",
				Usage:       "集群名称",
			})
		}
	}

	flags = append(flags, cli.StringFlag{
		Name:        "trace",
		Destination: &m.Trace,
		EnvVar:      "hydra_trace",
		Usage:       fmt.Sprintf("性能跟踪%v", supportTraces),
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "remote-logger,l",
		Destination: &m.remoteLogger,
		EnvVar:      "hydra_rpclog",
		Usage:       "启用远程日志",
	})
	flags = append(flags, cli.BoolFlag{
		Name:        "rqs,R",
		Destination: &m.RemoteQueryService,
		EnvVar:      "hydra-rqs",
		Usage:       "启用远程查询服务",
	})
	return flags
}
