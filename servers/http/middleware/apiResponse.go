package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"

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

		tp, content, err := nctx.Response.GetJSONRenderContent()
		writeTrace(getTrace(conf), tp, ctx, content)
		if err != nil && err.Error() != "" {
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
		case context.CT_PLAIN, context.CT_HTML, context.CT_OTHER:
			if v, ok := content.([]byte); ok {
				ctx.Data(nctx.Response.GetStatus(), tpName, v)
				return
			}
			ctx.Data(nctx.Response.GetStatus(), tpName, ([]byte)(content.(string)))
		default:
			ctx.JSON(nctx.Response.GetStatus(), content)
		}
	}
}
func writeTrace(b bool, tp int, ctx *gin.Context, c interface{}) {
	if !b {
		return
	}
	switch v := c.(type) {
	case []byte:
		setResponseRaw(ctx, string(v))
	case string:
		setResponseRaw(ctx, v)
	default:
		var buff = bytes.NewBufferString("")
		switch tp {
		case context.CT_XML:
			xml.NewEncoder(buff).Encode(c)
		default:
			json.NewEncoder(buff).Encode(c)
		}
		setResponseRaw(ctx, strings.Trim(buff.String(), "\n"))
		buff.Reset()
	}
}
