package middleware

import "github.com/gin-gonic/gin"

type HandlerFunc func(ctx *gin.Context)

func (h HandlerFunc) Handle(ctx *gin.Context) {
	h(ctx)
}
