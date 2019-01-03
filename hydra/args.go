package hydra

import (
	"github.com/urfave/cli"
)

type ArgContext struct {
	*cli.Context
	cmds     []cli.Flag
	Validate func() error
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
}
func (a *ArgContext) Append(c cli.Flag) {
	a.cmds = append(a.cmds, c)
}
