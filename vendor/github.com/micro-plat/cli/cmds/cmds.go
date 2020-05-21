package cmds

import (
	"github.com/urfave/cli"
)

var cmds []cli.Command = make([]cli.Command, 0, 4)

//Register 注册子
func Register(c ...cli.Command) {
	cmds = append(cmds, c...)
}

//GetCmds 获取所有命令
func GetCmds() []cli.Command {
	return cmds
}
