package rollback

import (
	"github.com/urfave/cli"
)

var backupFile string

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "f,file",
			Destination: &backupFile,
			Usage:       `-备份的本地服务配置文件`,
		},
	}
	return flags
}
