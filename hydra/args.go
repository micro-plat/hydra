package hydra

import (
	"github.com/urfave/cli"
)

type ArgContext struct {
	*cli.Context
	cmds     []cli.Flag
	Validate func() error
	RunMode  int
}

func newArgsContext() *ArgContext {
	return &ArgContext{
		cmds: make([]cli.Flag, 0, 0),
		Validate: func() error {
			return nil
		},
	}
}
func (a *ArgContext) setCtx(c *cli.Context) {
	a.Context = c
	switch c.Command.Name {
	case "install":
		a.RunMode = ModeInstall
	case "run":
		a.RunMode = ModeRun
	case "start":
		a.RunMode = ModeStart
	case "stop":
		a.RunMode = ModeStop
	case "status":
		a.RunMode = ModeStatus
	case "conf":
		a.RunMode = ModeConf
	case "remove":
		a.RunMode = ModeRemove
	case "registry":
		a.RunMode = ModeRegistry
	case "service":
		a.RunMode = ModeService

	}
}
func (a *ArgContext) Append(c cli.Flag) {
	a.cmds = append(a.cmds, c)
}
