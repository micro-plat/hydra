package main

import (
	"github.com/micro-plat/hydra/context"
)

func (rpc *rpcserver) handling() {
	rpc.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {

		return nil
	})
}
