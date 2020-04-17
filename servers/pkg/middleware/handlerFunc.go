package middleware

import "github.com/micro-plat/hydra/servers/pkg/dispatcher"

type HandlerFunc func(ctx *dispatcher.Context)

func (h HandlerFunc) Handle(ctx *dispatcher.Context) {
	h(ctx)
}
