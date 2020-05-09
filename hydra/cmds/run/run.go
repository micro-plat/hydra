package run

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/servers"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "run",
			Usage:  "运行服务",
			Flags:  getFlags(),
			Action: doRun,
		})
}

//doRun 服务启动
func doRun(c *cli.Context) (err error) {

	//1. 绑定应用程序参数
	if err := application.DefApp.Bind(); err != nil {
		logs.Log.Error(err)
		cli.ShowCommandHelp(c, c.Command.Name)
		return nil
		return nil
	}

	//2.创建trace性能跟踪
	if err = startTrace(application.Current().GetTrace()); err != nil {
		return
	}

	//3. 创建服务器
	server := servers.NewRspServers(application.Current().GetRegistryAddr(),
		application.Current().GetPlatName(), application.Current().GetSysName(),
		application.Current().GetServerTypes(), application.Current().GetClusterName())
	if err := server.Start(); err != nil {
		return err
	}

	//4. 堵塞当前进程，直到用户退出
	interrupt := make(chan os.Signal, 4)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM) //, syscall.SIGUSR1) //9:kill/SIGKILL,15:SIGTEM,20,SIGTOP 2:interrupt/syscall.SIGINT
LOOP:
	for {
		select {
		case <-time.After(time.Second * 120):
			debug.FreeOSMemory()
		case <-interrupt:
			break LOOP
		}
	}

	//5. 关闭服务器释放所有资源
	application.DefApp.Log().Info(application.AppName, "正在退出...")
	server.Shutdown()
	application.Current().Close()
	application.DefApp.Log().Info(application.AppName, "已安全退出")
	return nil

}
