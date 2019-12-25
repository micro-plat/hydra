package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/snappy"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/encoding/base64"
	"github.com/micro-plat/lib4go/logger"
)

func getUUID(c *dispatcher.Context) string {
	sid, ok := c.Request.GetHeader()["X-Request-Id"]
	if !ok || sid == "" {
		return logger.CreateSession()
	}
	return sid
}
func setUUID(c *dispatcher.Context, id string) {
	c.Request.GetHeader()["X-Request-Id"] = id
}

func setStartTime(c *dispatcher.Context) {
	c.Set("__start_time_", time.Now())
}
func setLogger(c *dispatcher.Context, l *logger.Logger) {
	c.Set("__logger_", l)
}
func getLogger(c *dispatcher.Context) *logger.Logger {
	l, _ := c.Get("__logger_")
	return l.(*logger.Logger)
}
func getExpendTime(c *dispatcher.Context) time.Duration {
	start, _ := c.Get("__start_time_")
	return time.Since(start.(time.Time))
}
func getJWTRaw(c *dispatcher.Context) interface{} {
	jwt, _ := c.Get("__jwt_")
	return jwt
}
func setJWTRaw(c *dispatcher.Context, v interface{}) {
	c.Set("__jwt_", v)
}
func getServiceName(c *dispatcher.Context) string {
	if service, ok := c.Get("__service_"); ok {
		return service.(string)
	}
	return ""
}
func setServiceName(c *dispatcher.Context, v string) {
	c.Set("__service_", v)
}
func setCTX(c *dispatcher.Context, r *context.Context) {
	c.Set("__context_", r)
}
func getTrace(cnf *conf.MetadataConf) bool {
	return cnf.GetMetadata("show-trace").(bool)
}
func getCTX(c *dispatcher.Context) *context.Context {
	result, _ := c.Get("__context_")
	if result == nil {
		return nil
	}
	return result.(*context.Context)
}
func getExt(c *dispatcher.Context) string {
	ext := make([]string, 0, 1)
	if f, ok := c.Get("__ext_param_name_"); ok {
		ext = append(ext, f.(string))
	}
	if v, ok := c.Get("__auth_tag_"); ok {
		ext = append(ext, v.(string))
	}
	return strings.Join(ext, " ")
}
func setResponseRaw(c *dispatcher.Context, raw string) {
	c.Set("__response_raw_", raw)
}
func getResponseRaw(c *dispatcher.Context) (string, bool) {
	if v, ok := c.Get("__response_raw_"); ok {
		return v.(string), true
	}
	return "", false
}
func setAuthTag(c *dispatcher.Context, ctx *context.Context) {
	if tag, ok := ctx.Response.GetParams()["__auth_tag_"]; ok {
		c.Set("__auth_tag_", tag)
	}

}

//ContextHandler api请求处理程序
func ContextHandler(exhandler interface{}, name string, engine string, service string, mSetting map[string]string, ext map[string]interface{}) dispatcher.HandlerFunc {

	handler, ok := exhandler.(servers.IExecuter)
	if !ok {
		panic("不是有效的servers.IExecuter接口")
	}

	return func(c *dispatcher.Context) {
		//处理输入参数
		ctx := context.GetContext(exhandler, name, engine, service, exhandler.(context.IContainer), makeQueyStringData(c), makeFormData(c), makeParamsData(c), makeSettingData(c, mSetting), makeExtData(c, ext), getLogger(c))

		defer setServiceName(c, service)
		defer setCTX(c, ctx)
		//调用执行引擎进行逻辑处理
		result := handler.Execute(ctx)
		if result != nil {
			ctx.Response.ShouldContent(result)
		}
		//处理错误err,5xx
		if err := ctx.Response.GetError(); err != nil {
			err = fmt.Errorf("error:%v", err)
			getLogger(c).Error(err)
			if !servers.IsDebug {
				err = errors.New("error:Internal Server Error")
			}
			return
		}
		//处理跳转3xx
		if url, ok := ctx.Response.IsRedirect(); ok {
			c.Redirect(ctx.Response.GetStatus(), url)
			return
		}

	}
}

func makeFormData(ctx *dispatcher.Context) IInputData {
	return newInputData(ctx.Request.GetForm, ctx.PostForm)
}
func makeQueyStringData(ctx *dispatcher.Context) IInputData {
	var p ParamData = make(map[string]string)
	return p
}
func makeParamsData(ctx *dispatcher.Context) IHeaderData {
	return newHeaderData(ctx.Params, ctx.Params.Get)
}
func makeSettingData(ctx *dispatcher.Context, m map[string]string) ParamData {
	return m
}

func makeExtData(c *dispatcher.Context, ext map[string]interface{}) map[string]interface{} {
	input := make(map[string]interface{})
	for k, v := range ext {
		if strings.HasPrefix(k, "__") {
			input[k] = v
			continue
		}
		input[fmt.Sprintf("__%s_", k)] = v
	}
	input["X-Request-Id"] = getUUID(c)
	input["__method_"] = strings.ToLower(c.Request.GetMethod())
	input["__header_"] = c.Request.GetHeader()
	input["__jwt_"] = getJWTRaw(c)
	input["__func_http_request_"] = c.Request
	input["__func_http_response_"] = c.Writer
	input["__binding_"] = func(obj interface{}) error {
		form := c.Request.GetForm()
		var buffer []byte
		var err error
		if body, ok := form["__body_"].(string); ok && len(form) == 1 {
			buffer = []byte(body)
		} else {
			buffer, err = json.Marshal(form)
			if err != nil {
				return err
			}
		}

		d := json.NewDecoder(bytes.NewBuffer(buffer))
		d.UseNumber()
		return d.Decode(obj)
	}
	input["__binding_with_"] = func(v interface{}, ct string) error {
		return input["__binding_"].(func(obj interface{}) error)(v)
	}
	input["__func_body_get_"] = func(ch string) (string, error) {
		if s, ok := c.Request.GetForm()["__body_"]; ok {
			if v, ok := c.Request.GetHeader()["__encode_snappy_"]; ok && v == "true" {
				buff, err := base64.DecodeBytes(s.(string))
				if err != nil {
					return "", fmt.Errorf("snappy压缩过的串必须是base64编码")
				}
				nbuffer, err := snappy.Decode(nil, buff)
				if err != nil {
					return "", fmt.Errorf("snappy.decode.err:%v %s", err, v)
				}
				return string(nbuffer), nil
			}
			return s.(string), nil
		}
		return "", errors.New("body读取错误")

	}
	input["__get_request_values_"] = func() map[string]interface{} {
		return c.Request.GetForm()
	}
	return input
}

//InputData 输入参数
type IInputData interface {
	Get(key string) (interface{}, bool)
	Keys() []string
}
type InputData struct {
	keys []string
	get  func(string) interface{}
}

//NewInputData 创建input data
func newInputData(v interface{}, get func(string) interface{}) *InputData {
	input := &InputData{
		get: get,
	}
	switch tp := v.(type) {
	case map[string][]string:
		input.keys = make([]string, 0, len(tp))
		for k := range tp {
			input.keys = append(input.keys, k)
		}
	case []dispatcher.Param:
		input.keys = make([]string, 0, len(tp))
		for _, k := range tp {
			input.keys = append(input.keys, k.Key)
		}
	}
	return input
}

//Get 获取指定键对应的数据
func (i InputData) Get(key string) (interface{}, bool) {
	return i.get(key), true
}

//Keys 获取所有KEY
func (i InputData) Keys() []string {
	return i.keys
}

//InputData 输入参数
type IHeaderData interface {
	Get(key string) (interface{}, bool)
	Keys() []string
}
type HeaderData struct {
	keys []string
	get  func(string) (string, bool)
}

//NewInputData 创建input data
func newHeaderData(tp []dispatcher.Param, get func(string) (string, bool)) *HeaderData {
	input := &HeaderData{
		get: get,
	}
	input.keys = make([]string, 0, len(tp))
	for _, k := range tp {
		input.keys = append(input.keys, k.Key)
	}
	return input
}

//Get 获取指定键对应的数据
func (i HeaderData) Get(key string) (interface{}, bool) {
	return i.get(key)
}

//Keys 获取所有KEY
func (i HeaderData) Keys() []string {
	return i.keys
}

//ParamData map参数数据
type ParamData map[string]string

//Get 获取指定键对应的数据
func (i ParamData) Get(key string) (interface{}, bool) {
	r, ok := i[key]
	return r, ok
}

//Keys 获取指定键对应的数据
func (i ParamData) Keys() []string {
	keys := make([]string, 0, len(i))
	for k := range i {
		keys = append(keys, k)
	}
	return keys
}
