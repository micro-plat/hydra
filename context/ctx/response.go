package ctx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
	"gopkg.in/yaml.v2"
)

var _ context.IResponse = &response{}

type rspns struct {
	status      int
	contentType string
	content     interface{}
}

type response struct {
	ctx         context.IInnerContext
	conf        server.IServerConf
	path        *rpath
	raw         rspns
	final       rspns
	noneedWrite bool
	log         logger.ILogger
	asyncWrite  func() error
	specials    []string
}

func NewResponse(ctx context.IInnerContext, conf server.IServerConf, log logger.ILogger, meta conf.IMeta) *response {
	return &response{
		ctx:   ctx,
		conf:  conf,
		path:  NewRpath(ctx, conf, meta),
		log:   log,
		final: rspns{status: http.StatusNotFound},
	}
}

//Header 设置头信息到response里
func (c *response) Header(k string, v string) {
	c.ctx.Header(k, v)
}

//ContentType 设置contentType
func (c *response) ContentType(v string) {
	c.ctx.Header("Content-Type", v)
}

//Abort 根据错误码与错误消息终止应用
func (c *response) Abort(s int, err error) {
	c.Write(s, err)
	c.ctx.Abort()
}

//Stop 停止当前服务执行
func (c *response) Stop(s int) {
	c.noneedWrite = true
	c.final.status = s
	c.ctx.Abort()
}

//StatusCode 设置response状态码
func (c *response) StatusCode(s int) {
	c.raw.status = s
	c.final.status = s
	c.ctx.WStatus(s)
}

//File 输入文件
func (c *response) File(path string) {
	if c.ctx.Written() {
		panic(fmt.Sprint("不能重复写入到响应流::", path))
	}
	c.ctx.File(path)
	c.ctx.Abort()
}

//WriteAny 将结果写入响应流，并自动处理响应码
func (c *response) WriteAny(v interface{}) error {
	return c.Write(http.StatusOK, v)
}

//NoNeedWrite 无需写入响应数据到缓存
func (c *response) NoNeedWrite(status int) {
	c.noneedWrite = true
	c.final.status = status
}

//Write 将结果写入响应流，自动检查内容处理状态码
func (c *response) Write(status int, content interface{}) error {
	if c.ctx.Written() {
		panic(fmt.Sprintf("不能重复写入到响应流:status:%d 已写入状态:%d", status, c.final.status))
	}

	//保存初始状态与结果
	c.raw.status, c.raw.content = status, content
	var ncontent interface{}

	//检查内容获取匹配的状态码
	c.final.status, ncontent = c.swapBytp(status, content)

	//检查内容类型并转换成字符串
	c.final.contentType, c.final.content = c.swapByctp(ncontent)

	//将编码设置到content type
	c.final.contentType = fmt.Sprintf(c.final.contentType, c.path.GetRouter().GetEncoding())

	//记录为原始状态
	c.raw.contentType = c.final.contentType

	//将写入操作处理为异步流程
	c.asyncWrite = func() error {
		return c.writeNow(c.final.status, c.final.contentType, c.final.content.(string))
	}
	return nil
}

//Render 修改实际渲染的内容
func (c *response) Render(status int, content string, ctp string) {
	if status != 0 {
		c.final.status = status
	}
	c.final.contentType = types.GetString(ctp, c.final.contentType)
	c.final.content = content
}
func (c *response) swapBytp(status int, content interface{}) (rs int, rc interface{}) {
	rs = status
	rc = content
	if content == nil {
		rc = ""
	}
	if status == 0 {
		rs = http.StatusOK
	}
	switch v := content.(type) {
	case errs.IError:
		rs = v.GetCode()
		rc = v.GetError().Error()
		c.log.Error(content)
		if global.IsDebug {
			rc = "Internal Server Error"
		}
	case error:
		if status >= http.StatusOK && status < http.StatusBadRequest {
			rs = http.StatusBadRequest
		}
		c.log.Error(content)
		rc = v.Error()
		if global.IsDebug {
			rc = "Internal Server Error"
		}
	default:
		return rs, rc
	}
	return rs, rc
}

func (c *response) swapByctp(content interface{}) (string, string) {
	ctp := c.getContentType()
	switch {
	case strings.Contains(ctp, "plain"):
		return ctp, fmt.Sprint(content)
	default:
		if content == nil || content == "" {
			return types.GetString(ctp, context.PLAINF), ""
		}
		tp := reflect.TypeOf(content).Kind()
		value := reflect.ValueOf(content)
		if tp == reflect.Ptr {
			value = value.Elem()
		}
		switch tp {
		case reflect.String:
			text := []byte(fmt.Sprint(content))
			switch {
			case (ctp == "" || strings.Contains(ctp, "json")) && json.Valid(text) && (bytes.HasPrefix(text, []byte("{")) ||
				bytes.HasPrefix(text, []byte("["))):
				return context.JSONF, content.(string)
			case (ctp == "" || strings.Contains(ctp, "xml")) && bytes.HasPrefix(text, []byte("<?xml")):
				return context.XMLF, content.(string)
			case strings.Contains(ctp, "html") && bytes.HasPrefix(text, []byte("<!DOCTYPE html")):
				return context.HTMLF, content.(string)
			case strings.Contains(ctp, "yaml"):
				return context.YAMLF, content.(string)
			case ctp == "" || strings.Contains(ctp, "plain"):
				return context.PLAINF, content.(string)
			default:
				return ctp, c.getString(ctp, map[string]interface{}{
					"data": content,
				})
			}
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			if ctp == "" {
				return context.PLAINF, fmt.Sprint(content)
			}
			return ctp, c.getString(ctp, map[string]interface{}{
				"data": content,
			})
		default:
			if ctp == "" {
				c.ContentType("application/json; charset=UTF-8")
				return context.JSONF, c.getString(context.JSONF, content)
			}
			return ctp, c.getString(ctp, content)
		}

	}
}
func (c *response) getContentType() string {
	if ctp := c.ctx.WHeader("Content-Type"); ctp != "" {
		return ctp
	}
	if ct, ok := c.conf.GetHeaderConf()["Content-Type"]; ok && ct != "" {
		return ct
	}
	return ""
}

//writeNow 将状态码、内容写入到响应流中
func (c *response) writeNow(status int, ctyp string, content string) error {
	if status >= http.StatusMultipleChoices && status < http.StatusBadRequest {
		c.ctx.Redirect(status, content)
		return nil
	}

	if c.path.GetRouter().IsUTF8() {
		c.ctx.Data(status, ctyp, []byte(content))
		return nil
	}
	buff, err := encoding.Encode(content, c.path.GetRouter().GetEncoding())
	if err != nil {
		return fmt.Errorf("输出时进行%s编码转换错误：%w %s", c.path.GetRouter().GetEncoding(), err, content)
	}
	c.ctx.Data(status, ctyp, buff)
	return nil
}

//Redirect 转跳g刚才gc
func (c *response) Redirect(code int, url string) {
	c.ctx.Redirect(code, url)
}

//AddSpecial 添加响应的特殊字符
func (c *response) AddSpecial(t string) {
	if c.specials == nil {
		c.specials = make([]string, 0, 1)
	}
	c.specials = append(c.specials, t)
}

//GetSpecials 获取多个响应特殊字符
func (c *response) GetSpecials() string {
	return strings.Join(c.specials, "|")
}

//GetRaw 获取原始响应请求
func (c *response) GetRaw() interface{} {
	return c.raw
}

//GetRawResponse 获取响应内容信息
func (c *response) GetRawResponse() (int, interface{}) {
	return c.raw.status, c.raw.content
}

//GetFinalResponse 获取响应内容信息
func (c *response) GetFinalResponse() (int, string) {
	if c.final.content == nil {
		return c.final.status, ""
	}
	return c.final.status, c.final.content.(string)
}

func (c *response) Flush() {
	if c.noneedWrite || c.asyncWrite == nil {
		return
	}
	if err := c.asyncWrite(); err != nil {
		panic(err)
	}

}
func (c *response) getString(ctp string, v interface{}) string {
	switch {
	case strings.Contains(ctp, "xml"):
		buff, err := xml.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(buff)
	case strings.Contains(ctp, "yaml"):
		buff, err := yaml.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(buff)
	case strings.Contains(ctp, "json"):
		buff, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(buff)
	}
	return fmt.Sprint(v)
}

func (c *response) getContent() string {
	return c.final.content.(string)
}
func (c *response) getStatus() int {
	return c.raw.status
}
