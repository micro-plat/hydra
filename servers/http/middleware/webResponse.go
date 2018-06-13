package middleware

import (
	"fmt"
	"path/filepath"

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
		switch nctx.Response.GetContentType() {
		case context.CT_JSON:
			ctx.JSON(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		case context.CT_XML:
			ctx.XML(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		case context.CT_YMAL:
			ctx.YAML(nctx.Response.GetStatus(), getMessage(nctx.Response.GetContent()))
		case context.CT_PLAIN:
			ctx.Data(nctx.Response.GetStatus(), "text/plain", []byte(fmt.Sprint(nctx.Response.GetContent())))
		case context.CT_HTML:
			ctx.Data(nctx.Response.GetStatus(), "text/html", []byte(fmt.Sprint(nctx.Response.GetContent())))
		default:
			if renderHTML(ctx, nctx.Response, conf) {
				return
			}
			ctx.Data(nctx.Response.GetStatus(), "text/plain", []byte(fmt.Sprint(nctx.Response.GetContent())))
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
