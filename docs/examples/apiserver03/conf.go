package main

import (
	"github.com/micro-plat/hydra"
)

func init() {
	hydra.OnReady(func() {
		/*hydra.Conf.API(":8080", api.WithTrace()).
			// APIKEY("sdfefefefefe").
			WhiteList(whitelist.NewIPList("/**", whitelist.WithIP("192.168.4.121"))).
			BlackList(blacklist.WithIP("192.168.4.120")).
			// Basic(basic.WithUP("admin", "123456")).
			// Fsa(fsa.CreateSecret(), fsa.WithInclude("/order/*")).
			// Jwt(jwt.WithExcludes("/member/**"), jwt.WithHeader()).
			// Static(static.WithArchive("./static.zip")).
			// Render(render.WithTmplt("/**", `{"id":{{get_status}}}`,
			// render.WithStatus(`{{get_req "id"}}`), render.WithContentType("")),
			// render.WithTmplt("/order/request", `{{get_status}}`)).
			Header()

		hydra.Conf.Vars().DB("db", oracle.New("hydra/hydra")).Queue("queue", lmq.New())
		*/
		hydra.Conf.API(":8083")
	})
}
