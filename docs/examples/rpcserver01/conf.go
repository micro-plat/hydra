package main

import (
	"github.com/micro-plat/hydra"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.RPC(":8092")
	})
}
