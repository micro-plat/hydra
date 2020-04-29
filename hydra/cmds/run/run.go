package run

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/servers"
	"github.com/urfave/cli"
	"github.com/zkfy/log"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "run",
			Usage:  "运行服务。前台运行，日志直接输出到客户端，输入ctl+c命令时退出服务",
			Flags:  getFlags(),
			Action: doRun,
		})
}

//doRun 服务启动
func doRun(c *cli.Context) (err error) {

	//1. 绑定应用程序参数
	if err := application.Bind(); err != nil {
		return err
	}

	//2. 创建服务器
	server := servers.NewRspServers(application.RegistryAddr, application.PlatName, application.SysName, application.ServerTypes, application.ClusterName)
	if err := server.Start(); err != nil {
		return err
	}

	//3. 堵塞当前进程，等服务结束
	signal.Notify(h.interrupt, os.Interrupt, os.Kill, syscall.SIGTERM) //, syscall.SIGUSR1) //9:kill/SIGKILL,15:SIGTEM,20,SIGTOP 2:interrupt/syscall.SIGINT
LOOP:
	for {
		select {
		case <-time.After(time.Second * 120):
			debug.FreeOSMemory()
		case <-h.interrupt:
			break LOOP
		}
	}

	//4. 关闭服务器释放所有资源
	log.Info(application.AppName, "正在退出...")
	server.Shutdown()
	log.Info(application.AppName, "已安全退出")
	return nil

}
