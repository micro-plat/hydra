package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type Handlers []Handler

func (c Handlers) DispFunc() []dispatcher.HandlerFunc {
	list := make([]dispatcher.HandlerFunc, 0, len(c))
	for _, item := range c {
		list = append(list, item.DispFunc())
	}
	return list
}

func (c Handlers) GinFunc() []gin.HandlerFunc {
	list := make([]gin.HandlerFunc, 0, len(c))
	for _, item := range c {
		list = append(list, item.GinFunc())
	}
	return list
}
func (c Handlers) Use(handler Handler) {
	c = append(c, handler)
}
