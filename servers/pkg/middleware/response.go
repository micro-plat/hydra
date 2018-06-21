package middleware

import (
	"errors"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
)

//Response 处理api返回值
func Response(conf *conf.MetadataConf) dispatcher.HandlerFunc {
	return func(ctx *dispatcher.Context) {
		ctx.Next()
		nctx := getCTX(ctx)
		if nctx == nil {
			return
		}
		defer nctx.Close()
		if err := nctx.Response.GetError(); err != nil {
			getLogger(ctx).Errorf("err:%v", err)
			if !servers.IsDebug {
				nctx.Response.ShouldContent(errors.New("请求发生错误"))
			}
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
		// switch nctx.Response.GetContentType() {
		// case 1:
		// 	ctx.SecureJSON(nctx.Response.GetStatus(), nctx.Response.GetContent())
		// case 2:
		// 	ctx.XML(nctx.Response.GetStatus(), nctx.Response.GetContent())
		// default:
		// 	if content, ok := nctx.Response.GetContent().(string); ok {
		// 		if (strings.HasPrefix(content, "[") || strings.HasPrefix(content, "{")) &&
		// 			(strings.HasSuffix(content, "}") || strings.HasSuffix(content, "]")) {
		// 			ctx.SecureJSON(nctx.Response.GetStatus(), nctx.Response.GetContent())
		// 		} else {
		// 			ctx.Data(nctx.Response.GetStatus(), "text/plain", []byte(nctx.Response.GetContent().(string)))
		// 		}
		// 		return
		// 	}
		// 	ctx.Data(nctx.Response.GetStatus(), "text/plain", []byte(fmt.Sprint(nctx.Response.GetContent())))
		// }
	}
}
