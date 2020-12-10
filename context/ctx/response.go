package ctx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
	"gopkg.in/yaml.v2"
)

var _ context.IResponse = &response{}

type rawrspns struct {
	status      int
	contentType string
	content     interface{}
}

type rspns struct {
	status      int
	contentType string
	content     string
}

type response struct {
	ctx         context.IInnerContext
	headers     types.XMap
	conf        app.IAPPConf
	path        *rpath
	raw         rawrspns
	final       rspns
	hasWrite    bool
	noneedWrite bool
	log         logger.ILogger
	asyncWrite  func() error
	specials    []string
}

//NewResponse 构建响应信息
func NewResponse(ctx context.IInnerContext, conf app.IAPPConf, log logger.ILogger, meta conf.IMeta) *response {
	return &response{
		ctx:  ctx,
		conf: conf,
		path: NewRpath(ctx, conf, meta),
		log:  log,
	}
}

//Header 设置头信息
func (c *response) Header(k string, v string) {
	c.ctx.Header(k, v)
}

//Header 获取头信息
func (c *response) GetHeaders() types.XMap {
	if c.headers != nil {
		return c.headers
	}
	hds := c.ctx.WHeaders()
	c.headers = make(map[string]interface{})
	for k, v := range hds {
		c.headers[k] = strings.Join(v, ",")
	}
	return c.headers
}

//ContentType 设置contentType
func (c *response) ContentType(v string) {
	c.ctx.Header("Content-Type", v)
}

//Abort 设置状态码,内容到响应流,并终止应用
func (c *response) Abort(s int, content ...interface{}) {
	defer c.ctx.Abort()
	defer c.Flush()
	c.Write(s, content...)
}

//File 将文件写入到响应流,并终止应用
func (c *response) File(path string) {
	defer c.ctx.Abort()
	if c.noneedWrite || c.ctx.Written() {
		return
	}
	c.noneedWrite = true
	c.raw.status = http.StatusOK
	c.final.status = http.StatusOK
	c.ctx.WStatus(http.StatusOK)
	c.ctx.File(path)
}

//NoNeedWrite 无需写入响应数据到缓存
func (c *response) NoNeedWrite(status int) {
	c.noneedWrite = true
	c.final.status = status
}

//JSON json输出响应内容
func (c *response) JSON(code int, data interface{}) interface{} {
	return c.Data(code, fmt.Sprintf(context.JSONF, c.path.GetEncoding()), data)
}

//XML xml输出响应内容
func (c *response) XML(code int, data interface{}) interface{} {
	return c.Data(code, fmt.Sprintf(context.XMLF, c.path.GetEncoding()), data)
}

//YAML yaml输出响应内容
func (c *response) YAML(code int, data interface{}) interface{} {
	return c.Data(code, fmt.Sprintf(context.YAMLF, c.path.GetEncoding()), data)
}

//Data 根据内容匹配响应内容
func (c *response) Data(code int, contentType string, data interface{}) interface{} {
	c.ContentType(contentType)
	if err := c.Write(code, data); err != nil {
		return err
	}
	return c.final.content
}

//Any 根据内容匹配响应内容与状态码
func (c *response) Any(data interface{}) interface{} {
	if err := c.Write(200, data); err != nil {
		return err
	}
	return c.final.content
}

//Write 检查内容并处理状态码,数据未写入响应流
func (c *response) Write(status int, ct ...interface{}) error {
	if c.noneedWrite {
		return fmt.Errorf("不能重复写入到响应流:status:%d 已写入状态:%d", status, c.final.status)
	}

	//1. 处理content
	var content interface{}
	if len(ct) > 0 {
		content = ct[0]
	}

	//2. 修改当前结果状态码与内容
	var ncontent interface{}
	c.final.status, ncontent = c.swapBytp(status, content)
	c.final.contentType, c.final.content = c.swapByctp(ncontent)
	if strings.Contains(c.final.contentType, "%s") {
		c.final.contentType = fmt.Sprintf(c.final.contentType, c.path.GetEncoding())
	}

	if c.hasWrite {
		return nil
	}

	//3. 保存初始状态与结果
	c.raw.status, c.raw.content, c.hasWrite, c.raw.contentType = status, content, true, c.final.contentType
	c.asyncWrite = func() error {
		return c.writeNow(c.final.status, c.final.contentType, c.final.content)
	}
	return nil
}

func (c *response) getContentType() string {
	if ctp := c.ctx.WHeader("Content-Type"); ctp != "" {
		return ctp
	}
	headerObj, err := c.conf.GetHeaderConf()
	if err != nil {
		return ""
	}
	if ct, ok := headerObj["Content-Type"]; ok && ct != "" {
		return ct
	}
	return ""
}

func (c *response) swapBytp(status int, content interface{}) (rs int, rc interface{}) {
	//处理状态码与响应内容的默认
	rs, rc = types.DecodeInt(status, 0, http.StatusOK), content

	switch v := content.(type) {
	case errs.IError:
		c.log.Error(content)
		rs, rc = v.GetCode(), types.DecodeString(global.IsDebug, false, "Internal Server Error", v.GetError().Error())
	case error:
		c.log.Error(content)
		rc = types.DecodeString(global.IsDebug, false, "Internal Server Error", v.Error())
		if status >= http.StatusOK && status < http.StatusBadRequest {
			rs = http.StatusBadRequest
		}
	}
	if content == nil {
		rc = ""
	}
	return rs, rc
}

func (c *response) swapByctp(content interface{}) (string, string) {

	//根据content-type反射内容进行输出
	vtpKind := getTypeKind(content)
	if ctp := c.getContentType(); ctp != "" {
		return ctp, c.getStringByCP(ctp, vtpKind, content)
	}

	//根据content确定 content-type
	if vtpKind == reflect.String {
		text := fmt.Sprint(content)
		switch {
		case strings.HasPrefix(text, "<!DOCTYPE html"):
			return context.HTMLF, text
		case strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">"):
			_, errx := mxj.BeautifyXml([]byte(text), "", "")
			if errx != nil {
				return context.PLAINF, text
			}
			return context.XMLF, text
		case json.Valid([]byte(text)) && (strings.HasPrefix(text, "{") ||
			strings.HasPrefix(text, "[")):
			return context.JSONF, text
		default:
			return context.PLAINF, text
		}

	} else if vtpKind == reflect.Struct || vtpKind == reflect.Map || vtpKind == reflect.Slice || vtpKind == reflect.Array {
		return context.JSONF, c.getStringByCP(context.JSONF, vtpKind, content)
	}
	return context.PLAINF, c.getStringByCP(context.JSONF, vtpKind, content)

}

func (c *response) getStringByCP(ctp string, tpkind reflect.Kind, content interface{}) string {
	if tpkind != reflect.Map && tpkind != reflect.Struct && tpkind != reflect.Slice && tpkind != reflect.Array {
		return fmt.Sprint(content)
	}

	switch {
	case strings.Contains(ctp, "xml"):
		if tpkind == reflect.Slice || tpkind == reflect.Array {
			panic("转化为xml必须是struct或者map,内容格式不正确")
		}
		if m := c.toMap(content); m != nil {
			if str, err := m.Xml(); err != nil {
				panic(err)
			} else {
				return string(str)
			}
		}
		if buff, err := xml.Marshal(content); err != nil {
			panic(err)
		} else {
			return string(buff)
		}
	case strings.Contains(ctp, "yaml"):
		if buff, err := yaml.Marshal(content); err != nil {
			panic(err)
		} else {
			return string(buff)
		}

	case strings.Contains(ctp, "json"):
		if buff, err := json.Marshal(content); err != nil {
			panic(err)
		} else {
			return string(buff)
		}
	default:
		return fmt.Sprint(content)
	}
}

func (c *response) toMap(content interface{}) mxj.Map {
	v := reflect.ValueOf(content)
	r := mxj.Map{}
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			r[types.GetString(key)] = v.MapIndex(key)
		}
	}
	return r
}

//writeNow 将状态码、内容写入到响应流中
func (c *response) writeNow(status int, ctyp string, content string) error {
	//301 302 303 307 308 这个地方会强制跳转到content 的路径。
	if status == http.StatusMovedPermanently || status == http.StatusFound || status == http.StatusSeeOther ||
		status == http.StatusTemporaryRedirect || status == http.StatusPermanentRedirect {
		//从header里面获取的Location
		location := content
		if l := c.ctx.WHeader("Location"); l != "" {
			location = l
		}
		c.ctx.Redirect(status, location)
		c.noneedWrite = true
		return nil
	}

	buff := []byte(content)
	e := c.path.GetEncoding()
	if e != encoding.UTF8 {
		buff1, err := encoding.Encode(content, e)
		if err == nil {
			buff = buff1
		}
	}
	c.ContentType(ctyp)
	c.ctx.Data(status, ctyp, buff)
	return nil
}

//WriteAny 向响应流中写入内容,状态码根据内容进行判断(不会立即写入)
func (c *response) WriteAny(v interface{}) error {
	return c.Write(http.StatusOK, v)
}

//Redirect 转跳g刚才gc
func (c *response) Redirect(code int, url string) {
	c.ctx.Redirect(code, url)
	c.noneedWrite = true
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
	return c.raw.content
}

//GetRawResponse 获取响应内容信息
func (c *response) GetRawResponse() (int, interface{}, string) {
	return c.raw.status, c.raw.content, c.raw.contentType
}

//GetFinalResponse 获取响应内容信息
func (c *response) GetFinalResponse() (int, string, string) {
	return c.final.status, c.final.content, c.final.contentType
}

//Flush 调用异步写入将状态码、内容写入到响应流中
func (c *response) Flush() {
	if c.noneedWrite || c.asyncWrite == nil || c.ctx.Written() {
		if c.final.status == 0 { //处理外部框架直接写入到流中,且输出日志状态为0的问题
			c.final.status = c.ctx.Status()
		}
		return
	}
	if err := c.asyncWrite(); err != nil {
		panic(err)
	}
	c.noneedWrite = true
}

func getTypeKind(c interface{}) reflect.Kind {
	if c == nil {
		return reflect.String
	}
	value := reflect.ValueOf(c)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Kind()
}
