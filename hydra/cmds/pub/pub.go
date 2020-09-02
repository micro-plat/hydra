package pub

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/hydra/global"
	"github.com/urfave/cli"
)

func init() {
	cmds.RegisterFunc(func() cli.Command {
		return cli.Command{
			Name:   "pub",
			Usage:  "远程发布服务，将当前应用发布到远程服务器，并启动",
			Flags:  getFlags(),
			Action: doPub,
		}
	})
}

func doPub(c *cli.Context) (err error) {

	//1.检查是否有管理员权限
	global.Current().Log().Pause()
	if err = global.CheckPrivileges(); err != nil {
		return err
	}

	//2.绑定请求参数
	if err = client.Bind(c.Args().Get(1), global.AppName, pwd); err != nil {
		return err
	}

	//3.登录到远程服务器
	if err = client.Login(); err != nil {
		return err
	}

	//4.切换工作目录
	if err := client.GoWorkDir(); err != nil {
		return err
	}

	//5. 上传文件
	if err := client.UploadFile(); err != nil {
		return err
	}

	//6.上传脚本
	path, err := client.UploadScript()
	if err != nil {
		return err
	}

	//7. 执行脚本
	if err := client.ExecScript(path); err != nil {
		return err
	}
	return nil
}
