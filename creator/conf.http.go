package creator

import (
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/router"

	"github.com/micro-plat/hydra/services"
)

type httpBuilder struct {
	CustomerBuilder
	tp          string
	fnGetRouter func(string) *services.ORouter
}

//newHTTP 构建http生成器
func newHTTP(tp string, address string, f func(string) *services.ORouter, opts ...api.Option) *httpBuilder {
	b := &httpBuilder{tp: tp, fnGetRouter: f, CustomerBuilder: make(map[string]interface{})}
	b.CustomerBuilder[ServerMainNodeName] = api.New(address, opts...)
	return b
}

//Load 加载路由
func (b *httpBuilder) Load() {
	routers, err := b.fnGetRouter(b.tp).GetRouters()
	if err != nil {
		panic(err)
	}
	b.CustomerBuilder[router.TypeNodeName] = routers
	return
}
