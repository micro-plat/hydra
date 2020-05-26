package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/auth/fsa"
	"github.com/micro-plat/hydra/conf/server/header"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080").Fsa(fsa.CreateSecret()).
			Header(header.WithCrossDomain())
	})
}
