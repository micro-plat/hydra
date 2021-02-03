package status

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/urfave/cli"
)

var isFixed bool

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "status",
			Usage:  "查询状态，查询服务器运行、停止状态",
			Flags:  pkgs.GetFixedFlags(&isFixed),
			Action: doStatus,
		}
	})
}

func doStatus(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()
	//3.创建本地服务
	hydraSrv, err := pkgs.GetService(c, isFixed)
	if err != nil {
		return err
	}
	status, err := hydraSrv.Status()
	return pkgs.GetCmdsResult(hydraSrv.DisplayName, "Status", err, statusMap[status])
}

var statusMap = map[service.Status]string{
	service.StatusRunning: "Running",
	service.StatusStopped: "Stopped",
	service.StatusUnknown: "Unknown",
}
