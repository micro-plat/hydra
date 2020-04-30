package remove

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/hydra/cmds/daemon"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "remove",
			Usage:  "删除服务。应用启动参数发生变化后，需调用remove删除本地服务后再重新安装",
			Action: doRemove,
		})
}
func doRemove(c *cli.Context) (err error) {
	service, err := daemon.New(application.AppName, application.AppName)
	if err != nil {
		return err
	}
	msg, err := service.Remove()
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
