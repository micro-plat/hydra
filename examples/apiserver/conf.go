package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/registry/conf/server/header"
)

func init() {
	hydra.Conf.Ready(func() {
		hydra.Conf.API(":8080").Jwt(jwt.WithExclude("/order/request/2")).Header(header.WithCrossDomain(), header.WithHeader("k", "v"))
	})
}
