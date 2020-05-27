package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/static"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080").
			// Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
			Jwt(jwt.WithExclude("/member/**", "/order/**"), jwt.WithHeader()).
			Static(static.WithArchive("./static.zip")).
			Header(header.WithCrossDomain())
	})
}
