package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080", api.WithTrace()).
			// Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
			Jwt(jwt.WithExclude("/member/**"), jwt.WithHeader()).
			Static(static.WithArchive("./static.zip")).
			Render(render.WithTmplt("/member/*", `{"id":{{get_status}}}`), render.WithTmplt("/order/request", `success{{get_path}}`)).
			Header(header.WithCrossDomain())
	})
}
