package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
)

//APIResponse 处理api返回值
func APIResponse(conf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		nctx := getCTX(ctx)
		if nctx == nil {
			return
		}
		if _, ok := nctx.Response.IsRedirect(); ok {
			return
		}

		if ctx.Writer.Written() {
			return
		}

		switch nctx.Response.GetContentType() {
		case context.CT_XML:
			ctx.XML(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		case context.CT_YMAL:
			ctx.YAML(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		case context.CT_PLAIN:
			ctx.Data(nctx.Response.GetStatus(), "text/plain", []byte(fmt.Sprint(nctx.Response.GetContent())))
		case context.CT_HTML:
			ctx.Data(nctx.Response.GetStatus(), "text/html", []byte(fmt.Sprint(nctx.Response.GetContent())))
		default:
			ctx.JSON(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		}
	}
}
func getMessage(i interface{}) interface{} {
	switch v := i.(type) {
	case string:
		return json.RawMessage(v)
	case bool:
		return map[string]interface{}{
			"bool": v,
		}
	case int, float32, float64:
		return map[string]interface{}{
			"num": v,
		}
	case error:
		return map[string]interface{}{
			"err": v.Error(),
		}
	default:
		return i
	}
}
