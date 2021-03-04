package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

type node struct {
	path     string
	service  string
	actions  []string
	handlers middleware.Handlers
}

func (n *node) GetGinHandlers(tps ...string) gin.HandlersChain {
	handlersChain := make(gin.HandlersChain, len(n.handlers))
	for i := range n.handlers {
		handlersChain[i] = n.handlers[i].GinFunc(tps...)
	}
	return handlersChain
}

func (n *node) GetDispHandlers(tps ...string) dispatcher.HandlersChain {
	handlersChain := make(dispatcher.HandlersChain, len(n.handlers))
	for i := range n.handlers {
		handlersChain[i] = n.handlers[i].DispFunc(tps...)
	}
	return handlersChain
}
