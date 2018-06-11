package middleware

import (
	"errors"
	"fmt"

	"github.com/micro-plat/hydra/servers"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
)

//WSResponse 处理api返回值
func WSResponse(conf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		nctx := getCTX(ctx)
		if nctx == nil {
			return
		}
		defer nctx.Close()
		if err := nctx.Response.GetError(); err != nil {
			getLogger(ctx).Error(err)
			if !servers.IsDebug {
				nctx.Response.ShouldContent(errors.New("请求发生错误"))
			}
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
