package hydra

import "github.com/micro-plat/hydra/global"

type ucli struct {
	Run  global.IUCLI
	Conf global.IUCLI
}

func newUCli() *ucli {
	return &ucli{
		Run:  global.RunCli,
		Conf: global.ConfCli,
	}
}
