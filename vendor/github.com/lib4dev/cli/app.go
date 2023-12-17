package cli

import (
	"os"
	"path/filepath"

	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/urfave/cli"
)

//VERSION 版本号
var VERSION = "0.0.1"

//App  cli app
type App struct {
	*cli.App
	log *logs.Logger
	*option
}

//Start 启动应用程序
func (a *App) Start() {
	if err := a.Run(os.Args); err != nil {
		a.log.Error(err)
	}
}

//New 创建app
func New(opts ...Option) *App {

	app := &App{log: logs.New(), option: &option{version: VERSION, usage: "A new cli application"}}
	for _, opt := range opts {
		opt(app.option)
	}

	app.App = cli.NewApp()
	app.App.Name = filepath.Base(os.Args[0])
	app.App.Version = app.version
	app.App.Usage = app.usage
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help,h",
		Usage: "查看帮助信息",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,v",
		Usage: "查看版本信息",
	}
	app.App.Commands = cmds.GetCmds()
	return app
}
