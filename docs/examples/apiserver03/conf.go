package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/fsa"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080", api.WithTrace()).
			Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
			Jwt(jwt.WithExclude("/member/**"), jwt.WithHeader()).
			Static(static.WithArchive("./static.zip")).
			Render(render.WithTmplt("/**", `{"id":{{get_status}}}`,
				render.WithStatus(`{{get_req "id"}}`), render.WithContentType("")),
				render.WithTmplt("/order/request", `{{get_status}}`)).
			Header(header.WithCrossDomain())

		hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
	})
}
