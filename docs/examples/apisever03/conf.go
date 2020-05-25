package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/header"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080").
			Jwt(jwt.WithExclude("/**")).
			Header(header.WithCrossDomain())
	})
}
