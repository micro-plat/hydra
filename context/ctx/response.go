package ctx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
)

var _ context.IResponse = &response{}

type response struct {
	ctx      context.IInnerContext
	conf     server.IServerConf
	content  interface{}
	log      logger.ILogger
	specials []string
}

//Header 设置头信息到response里
func (c *response) SetHeader(k string, v string) {
	c.ctx.Header(k, v)
}

//ContentType 设置contentType
func (c *response) ContentType(v string) {
	c.ctx.Header("Content-Type", v)
}

//Abort 根据错误码终止应用
func (c *response) Abort(s int) {
	c.Write(s, nil)
	c.ctx.Abort()
}

//AbortWithError 根据错误码与错误消息终止应用
func (c *response) AbortWithError(s int, err error) {
	c.Write(s, err)
	c.ctx.Abort()
}

//GetStatusCode 获取response状态码
func (c *response) GetStatusCode() int {
	return c.ctx.Status()
}

//SetStatusCode 设置response状态码
func (c *response) SetStatusCode(s int) {
	c.ctx.WStatus(s)
}

//Written 响应是否已写入
func (c *response) Written() bool {
	return c.ctx.Written()
}

//File 输入文件
func (c *response) File(path string) {
	if c.ctx.Written() {
		panic(fmt.Sprint("不能重复写入到响应流::", path))
	}
	c.File(path)
}

//WriteAny 将结果写入响应流，并自动处理响应码
func (c *response) WriteAny(v interface{}) error {
	return c.Write(200, v)
}

//Write 将结果写入响应流，自动检查内容处理状态码
func (c *response) Write(status int, content interface{}) error {
	if c.ctx.Written() {
		panic(fmt.Sprint("不能重复写入到响应流:status:", status, content))
	}
	switch v := content.(type) {
	case errs.IError:
		status = v.GetCode()
		content = v.GetError().Error()
		c.log.Error(content)
		if global.IsDebug {
			content = "Internal Server Error"
		}
	case error:
		if status >= 200 && status < 400 {
			status = 400
		}
		content = v.Error()
		c.log.Error(content)
		if global.IsDebug {
			content = "Internal Server Error"
		}
	}
	return c.writeNow(status, c.swap(content))
}
func (c *response) swap(content interface{}) interface{} {
	ctp := c.getContentType()
	switch {
	case strings.Contains(ctp, "plain"):
		c.ContentType(ctp)
		return fmt.Sprint(content)
	default:
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
				c.ContentType("application/json; charset=UTF-8")
				return content
			case (ctp == "" || strings.Contains(ctp, "xml")) && bytes.HasPrefix(text, []byte("<?xml")):
				c.ContentType("application/xml; charset=UTF-8")
				return content
			case strings.Contains(ctp, "html") && bytes.HasPrefix(text, []byte("<!DOCTYPE html")):
				c.ContentType("text/html; charset=UTF-8")
				return content
			case strings.Contains(ctp, "yaml"):
				c.ContentType(ctp)
				return content
			case ctp == "" || strings.Contains(ctp, "plain"):
				c.ContentType("text/plain; charset=UTF-8")
				return content
			default:
				c.ContentType(ctp)
				return map[string]interface{}{
					"data": content,
				}
			}
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			if ctp == "" {
				c.ContentType("text/plain; charset=UTF-8")
				return content
			}
			return map[string]interface{}{
				"data": content,
			}
		default:
			if ctp == "" {
				c.ContentType("application/json; charset=UTF-8")
				return content
			}
			return content
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
func (c *response) writeNow(status int, content interface{}) error {
	if c.ctx.WHeader("Content-Type") == "" {
		c.ContentType("application/json; charset=UTF-8")
	}
	c.content = content
	tpName := c.ctx.WHeader("Content-Type")
	switch v := content.(type) {
	case []byte:
		c.ctx.Data(status, tpName, v)
		return nil
	case string:
		c.ctx.Data(status, tpName, []byte(v))
		return nil
	}
	switch {
	case strings.Contains(tpName, "xml"):
		c.ctx.XML(status, content)
	case strings.Contains(tpName, "yaml"):
		c.ctx.YAML(status, content)
	default:
		c.ctx.JSON(status, content)
	}

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

//GetResponse 获取响应内容信息
func (c *response) GetResponse() string {
	if c.content == nil {
		return "[nil]"
	}
	switch v := c.content.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	default:
		if buff, err := json.Marshal(c.content); err == nil {
			return string(buff)
		}
		return fmt.Sprint(c.content)
	}
}
