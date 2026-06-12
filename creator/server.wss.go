package creator

import (
	"github.com/micro-plat/hydra/conf/server/wss"
	"github.com/micro-plat/hydra/global"
)

type wssBuilder struct {
	BaseBuilder
	tp string
}

func newWSS(opts ...wss.Option) *wssBuilder {
	var side wss.Option
	if len(opts) > 0 {
		side = opts[0]
	}
	if side == nil {
		side = wss.WithServerSide()
	}
	b := &wssBuilder{tp: side.Type(), BaseBuilder: make(map[string]interface{})}
	switch b.tp {
	case global.WSSClient:
		s := wss.NewClientSide()
		side.Apply(s)
		b.BaseBuilder[ServerMainNodeName] = s
	default:
		b.tp = global.WSSServer
		s := wss.NewServerSide()
		side.Apply(s)
		b.BaseBuilder[ServerMainNodeName] = s
	}
	return b
}

func (b *wssBuilder) Load() {
}

func (b *wssBuilder) Routes(routes ...wss.Route) *wssBuilder {
	b.BaseBuilder[wss.RouteNodeName] = wss.NewRoutes(routes...)
	return b
}
