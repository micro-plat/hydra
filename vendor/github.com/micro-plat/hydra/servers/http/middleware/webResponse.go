package middleware

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"

	"github.com/gin-gonic/gin"
)

//WebResponse 处理web返回值
func WebResponse(conf *conf.MetadataConf) gin.HandlerFunc {
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

		tp, content, err := nctx.Response.GetHTMLRenderContent()
		writeTrace(getTrace(conf), tp, ctx, content)
		if err != nil && err.Error() != "" {
			getLogger(ctx).Error(err)
			ctx.JSON(nctx.Response.GetStatus(), map[string]interface{}{"err": err})
			return
		}
		tpName := context.ContentTypes[tp]
		switch tp {
		case context.CT_JSON:
			ctx.JSON(nctx.Response.GetStatus(), content)
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
		case context.CT_PLAIN:
			ctx.Data(nctx.Response.GetStatus(), tpName, content.([]byte))
		default:
			if renderHTML(ctx, nctx.Response, conf) {
				return
			}
			html := fmt.Sprint(content)
			if strings.HasPrefix(html, "<!DOCTYPE html") {
				ctx.Data(nctx.Response.GetStatus(), tpName, []byte(html))
				return
			}
			ctx.Data(nctx.Response.GetStatus(), context.ContentTypes[context.CT_PLAIN], []byte(html))
		}
	}
}
func renderHTML(ctx *gin.Context, response context.IResponse, cnf *conf.MetadataConf) bool {
	files, ok := cnf.GetMetadata("viewFiles").([]string)
	if !ok {
		return false
	}
	root := cnf.GetMetadata("view").(*conf.View).Path
	viewPath := filepath.Join(root, fmt.Sprintf("%s.html", getServiceName(ctx)))
	for _, f := range files {
		if f == viewPath {
			ctx.HTML(response.GetStatus(), filepath.Base(viewPath), response.GetContent())
			return true
		}
	}
	return false
}
