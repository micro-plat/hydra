package backup

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/daemon"
	"github.com/micro-plat/lib4go/errs"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "backup",
			Usage:  "备份服务",
			Flags:  getFlags(),
			Action: doBackup,
		}
	})
}

func doBackup(c *cli.Context) (err error) {

	//关闭日志显示
	global.Current().Log().Pause()
	service, err := daemon.New(global.Def.GetLongAppName(), global.Usage)
	if err != nil {
		return err
	}
	msg, err := service.Backup(backupFile)
	if err != nil {
		return err
	}
	return errs.NewIgnoreError(0, msg)
}
