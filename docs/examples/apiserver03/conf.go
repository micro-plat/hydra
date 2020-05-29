package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/render"
)

// Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
//Jwt(jwt.WithExclude("/member/**"), jwt.WithHeader()).
//Static(static.WithArchive("./static.zip")).
func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080", api.WithTrace()).
			Render(render.WithTmplt("/**", `{"id":{{get_status}}}`,
				render.WithStatus(`{{get_req "id"}}`), render.WithContentType("")),
				render.WithTmplt("/order/request", `{{get_status}}`)).
			Header(header.WithCrossDomain())
	})
}
