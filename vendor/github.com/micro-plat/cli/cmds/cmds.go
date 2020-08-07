package cmds

import (
	"sync"

	"github.com/urfave/cli"
)

var cmds []cli.Command = make([]cli.Command, 0, 4)
var funcs []func() cli.Command = make([]func() cli.Command, 0, 1)
var once sync.Once

//Register 注册子
func Register(c ...cli.Command) {
	cmds = append(cmds, c...)
}

//RegisterFunc 注册函数，用于异步加载
func RegisterFunc(f ...func() cli.Command) {
	funcs = append(funcs, f...)
}

//GetCmds 获取所有命令
func GetCmds() []cli.Command {
	once.Do(func() {
		for _, f := range funcs {
			Register(f())
		}
	})
	return cmds
}
