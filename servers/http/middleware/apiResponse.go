package middleware

import (
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
		if url, ok := nctx.Response.IsRedirect(); ok {
			ctx.Redirect(nctx.Response.GetStatus(), url)
			return
		}

		if ctx.Writer.Written() {
			return
		}

		tp, content, err := nctx.Response.GetRenderContent(context.CT_JSON)
		if err != nil {
			getLogger(ctx).Error(err)
			ctx.JSON(nctx.Response.GetStatus(), map[string]interface{}{"err": err})
			return
		}
		tpName := context.ContentTypes[tp]
		switch tp {
		case context.CT_XML:
			if v, ok := content.([]byte); ok {
				ctx.Data(nctx.Response.GetStatus(), tpName, v)
				return
			}
			ctx.XML(nctx.Response.GetStatus(), content)
		case context.CT_YMAL:
			if v, ok := content.([]byte); ok {
				ctx.Data(nctx.Response.GetStatus(), tpName, v)
				return
			}
			ctx.YAML(nctx.Response.GetStatus(), content)
		case context.CT_PLAIN, context.CT_HTML:
			ctx.Data(nctx.Response.GetStatus(), tpName, content.([]byte))
		default:
			ctx.JSON(nctx.Response.GetStatus(), content)
		}
	}
}
