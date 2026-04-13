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

	var nwriter interface{}
	switch strings.ToLower(ctx.GetType()) {
	case "gin":
		writer := ctx.GetWriter().(gin.ResponseWriter)
		nwriter = newGinWriter(writer, ctx, g.level)
		ctx.SetWriter(nwriter)
	default:
		writer := ctx.GetWriter().(dispatcher.ResponseWriter)
		nwriter = newDispWriter(writer, ctx, g.level)
		ctx.SetWriter(nwriter)
	}

	// 使用 defer 确保在请求结束时关闭 gzip writer
	// defer 会在 ctx.Next() 返回后执行，也就是所有中间件（包括 Flush）完成后
	defer func() {
		if closer, ok := nwriter.(interface{ Close() }); ok {
			closer.Close()
		}
	}()

	ctx.Next()
}
