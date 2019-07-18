package rgsts

import (
	"strings"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/urfave/cli"
)

var Rgst registry.IRegistry

//Log  日志组件
var Log logger.ILogging

var cmds []cli.Command
var Root string

//Register 注册子
func Register(c ...cli.Command) {
	cmds = append(cmds, c...)
}
func init() {
	cmds = make([]cli.Command, 0, 4)
	Log = newLogger()
}
func GetCmds() []cli.Command {
	return cmds
}
func GetRoot() string {
	if !strings.HasPrefix(Root, "/") {
		return "/" + Root
	}
	return Root
}
