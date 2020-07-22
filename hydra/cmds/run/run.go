package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/rlog"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "run",
			Usage:  "运行服务",
			Flags:  getFlags(),
			Action: doRun,
		}
	})
}

//doRun 服务启动
func doRun(c *cli.Context) (err error) {

	//1. 绑定应用程序参数
	if err := global.Def.Bind(c); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
	}

	//3. 注册远程日志组件
	if err = rlog.Registry(global.Def.PlatName, global.Def.RegistryAddr); err != nil {
		logs.Log.Error(err)
		return nil
	}

	//4.创建trace性能跟踪
	if err = startTrace(global.Current().GetTrace()); err != nil {
		return
	}

	//5. 处理本地内存作为注册中心的服务发布问题
	if registry.GetProto(global.Current().GetRegistryAddr()) == registry.LocalMemory {
		if err := pkgs.Pub2Registry(true); err != nil {
			return err
		}
	}

	//6. 创建服务器
	server := servers.NewRspServers(global.Current().GetRegistryAddr(),
		global.Current().GetPlatName(), global.Current().GetSysName(),
		global.Current().GetServerTypes(), global.Current().GetClusterName())
	if err := server.Start(); err != nil {
		return err
	}

	//7. 堵塞当前进程，直到用户退出
	interrupt := make(chan os.Signal, 4)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR2) //, syscall.SIGUSR1) //9:kill/SIGKILL,15:SIGTEM,20,SIGTOP 2:interrupt/syscall.SIGINT
	var signal os.Signal
LOOP:
	for {
		select {
		case signal = <-interrupt:
			break LOOP
		}
	}

	//8. 关闭服务器释放所有资源
	global.Def.Log().Info(global.AppName, fmt.Sprintf("正在退出[%v]...", signal))

	//关闭服务器
	server.Shutdown()

	//关闭各服务
	if err := services.Def.Close(); err != nil {
		global.Def.Log().Error("err:", err)
	}

	//通知关闭各组件
	global.Current().Close()

	global.Def.Log().Info(global.AppName, "已安全退出")
	return nil

}
