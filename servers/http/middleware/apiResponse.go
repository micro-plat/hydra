package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//APIResponse 处理api返回值
func APIResponse(xconf *conf.MetadataConf) gin.HandlerFunc {
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
		writeTrace(getTrace(xconf), tp, ctx, content)
		if err != nil && err.Error() != "" {
			getLogger(ctx).Error(err)
			tp = context.CT_JSON
			content = err.Error()
			nctx.Response.ShouldContent(err)
		}

		//检查是否配置输出模板
		resTemplate, ok := xconf.GetMetadata("__response_conf_").(*conf.Response)
		if !ok {
			write2Response(ctx, tp, nctx.Response.GetStatus(), content)
			return
		}
		ok, template := resTemplate.GetTemplate(nctx.Service)
		if !ok {
			write2Response(ctx, tp, nctx.Response.GetStatus(), content)
			return
		}

		//翻译模板进行输出
		status, content, err := getResponseContent(template, nctx, tp, content)
		if err != nil && err.Error() != "" {
			getLogger(ctx).Error(err)
			content = map[string]interface{}{"err": err}
		}
		write2Response(ctx, tp, status, content)

	}
}

func getResponseContent(c *conf.Template, ctx *context.Context, t int, sc interface{}) (int, interface{}, error) {
	status := ctx.Response.GetStatus()
	input := types.NewXMap()
	input.MergeMap(c.Params)
	input.MergeMap(ctx.Response.GetParams())
	input.SetValue("status", ctx.Response.GetStatus())
	input.SetValue("param", ctx.Request.Param.GetMaps())
	input.SetValue("querystring", ctx.Request.QueryString.GetMaps())
	input.SetValue("form", ctx.Request.Form.GetMaps())
	if err := ctx.Response.GetError(); err != nil {
		input.SetValue("err", err.Error())
	} else {
		input.SetValue("data", ctx.Response.GetContent())
	}
	if !input.Has("sdata") {
		input.SetValue("sdata", sc)
	}

	//翻译状态码
	code, err := c.GetStatus(input.ToMap())
	if err != nil {
		return status, nil, err
	}

	//翻译模块
	result, err := c.GetContent(input.ToMap())
	if err != nil {
		return status, nil, err
	}
	if v := types.GetInt(code, 0); v != 0 {
		status = v
	}
	if result != "" {
		sc = result
	}
	return status, sc, nil

}

//将指定的状态码，内容输出到响应流
func write2Response(ctx *gin.Context, tp int, status int, content interface{}) {
	tpName := context.ContentTypes[tp]
	switch v := content.(type) {
	case []byte:
		ctx.Data(status, tpName, v)
		return
	case string:
		ctx.Data(status, tpName, []byte(v))
		return
	}
	switch tp {
	case context.CT_XML:
		ctx.XML(status, content)
	case context.CT_YMAL:
		ctx.YAML(status, content)
	default:
		ctx.JSON(status, content)
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
