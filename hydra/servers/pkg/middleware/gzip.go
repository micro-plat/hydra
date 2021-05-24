package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

func Gzip(level int) Handler {
	return newGzipHandler(level).Handle
}

type gzipHandler struct {
	level int
}

func newGzipHandler(level int) *gzipHandler {
	return &gzipHandler{
		level: level,
	}
}

func (g *gzipHandler) Handle(ctx IMiddleContext) {
	processor, err := ctx.APPConf().GetProcessorConf()
	if err != nil {
		ctx.Response().Abort(http.StatusNotExtended, err)
		return
	}

	//未启用gzip
	if !processor.EnableGzip {
		ctx.Next()
		return
	}

	switch strings.ToLower(ctx.GetType()) {
	case "gin":
		writer := ctx.GetWriter().(gin.ResponseWriter)
		nwriter := newGinWriter(writer, ctx, g.level)
		ctx.SetWriter(nwriter)
		ctx.Response().OnFlush(nwriter.Close)
	default:
		writer := ctx.GetWriter().(dispatcher.ResponseWriter)
		nwriter := newDispWriter(writer, ctx, g.level)
		ctx.SetWriter(nwriter)
		ctx.Response().OnFlush(nwriter.Close)
	}
	ctx.Next()
}
