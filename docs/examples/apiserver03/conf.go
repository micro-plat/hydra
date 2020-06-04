package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
)

func init() {
	hydra.Conf.OnReady(func() {
		hydra.Conf.API(":8080", api.WithTrace()).
			Basic(basic.WithUP("admin", "123456")).
			// Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
			// Jwt(jwt.WithExcludes("/member/**"), jwt.WithHeader()).
			// Static(static.WithArchive("./static.zip")).
			// Render(render.WithTmplt("/**", `{"id":{{get_status}}}`,
			// render.WithStatus(`{{get_req "id"}}`), render.WithContentType("")),
			// render.WithTmplt("/order/request", `{{get_status}}`)).
			Header()

		hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
	})
}
