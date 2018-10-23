package middleware

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/logger"
)

func getUUID(c *gin.Context) string {
	if v, ok := c.Get("__hydra_sid_"); ok {
		return v.(string)
	}
	ck, err := c.Request.Cookie("hydra_sid")
	if err != nil || ck == nil || ck.Value == "" {
		return logger.CreateSession()
	}
	return ck.Value
}
func setUUID(c *gin.Context, id string) {
	c.Set("__hydra_sid_", id)
}
func setStartTime(c *gin.Context) {
	c.Set("__start_time_", time.Now())
}

func setExt(c *gin.Context, name string) {
	c.Set("__ext_param_name_", name)
}
func getExt(c *gin.Context) string {
	if f, ok := c.Get("__ext_param_name_"); ok {
		return f.(string)
	}
	return ""
}

func setLogger(c *gin.Context, l *logger.Logger) {
	c.Set("__logger_", l)
}
func getLogger(c *gin.Context) *logger.Logger {
	l, _ := c.Get("__logger_")
	return l.(*logger.Logger)
}
func getExpendTime(c *gin.Context) time.Duration {
	start, _ := c.Get("__start_time_")
	return time.Since(start.(time.Time))

}
func getJWTRaw(c *gin.Context) interface{} {
	jwt, _ := c.Get("__jwt_")
	return jwt
}
func setJWTRaw(c *gin.Context, v interface{}) {
	c.Set("__jwt_", v)
}

func getIsCircuitBreaker(c *gin.Context) bool {
	if b, ok := c.Get("__is_circuit_breaker_"); ok {
		return b.(bool)
	}
	return false
}
func setIsCircuitBreaker(c *gin.Context, v bool) {
	c.Set("__is_circuit_breaker_", v)
}

func getServiceName(c *gin.Context) string {
	if service, ok := c.Get("__service_"); ok {
		return service.(string)
	}
	return ""
}
func setServiceName(c *gin.Context, v string) {
	c.Set("__service_", v)
}
func setCTX(c *gin.Context, r *context.Context) {
	c.Set("__context_", r)
}
func getCTX(c *gin.Context) *context.Context {
	result, _ := c.Get("__context_")
	if result == nil {
		return nil
	}
	return result.(*context.Context)
}
func getTrace(cnf *conf.MetadataConf) bool {
	return cnf.GetMetadata("show-trace").(bool)
}
func setResponseRaw(c *gin.Context, raw string) {
	c.Set("__response_raw_", raw)
}
func getResponseRaw(c *gin.Context) (string, bool) {
	if v, ok := c.Get("__response_raw_"); ok {
		return v.(string), true
	}
	return "", false
}

func setMetadataConf(c *gin.Context, cnf *conf.MetadataConf) {
	c.Set("__metadata-conf_", cnf)
}
func getMetadataConf(c *gin.Context) *conf.MetadataConf {
	v, ok := c.Get("__metadata-conf_")
	if !ok {
		return nil
	}
	return v.(*conf.MetadataConf)
}

//ContextHandler api请求处理程序
func ContextHandler(exhandler interface{}, name string, engine string, service string, mSetting map[string]string) gin.HandlerFunc {
	handler, ok := exhandler.(servers.IExecuter)
	if !ok {
		panic("不是有效的servers.IExecuter接口")
	}

	return func(c *gin.Context) {
		//处理输入参数
		ctn, _ := exhandler.(context.IContainer)
		ctx := context.GetContext(exhandler, name, engine, service, ctn, makeQueyStringData(c), makeFormData(c), makeParamsData(c), makeSettingData(c, mSetting), makeExtData(c), getLogger(c))

		defer setServiceName(c, ctx.Service)
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
			ctx.Response.ShouldContent(err)
		}
	}
}

func makeFormData(ctx *gin.Context) InputData {
	if ctx.ContentType() == binding.MIMEPOSTForm ||
		ctx.ContentType() == binding.MIMEMultipartPOSTForm {
		ctx.Request.ParseForm()
		ctx.Request.ParseMultipartForm(32 << 20)
	}

	return ctx.GetPostForm
}
func makeQueyStringData(ctx *gin.Context) InputData {
	return ctx.GetQuery
}
func makeParamsData(ctx *gin.Context) InputData {
	return ctx.Params.Get
}

func makeMapData(m map[string]interface{}) MapData {
	return m
}

func makeSettingData(ctx *gin.Context, m map[string]string) ParamData {
	return m
}

func makeExtData(c *gin.Context) map[string]interface{} {
	input := make(map[string]interface{})
	input["__hydra_sid_"] = getUUID(c)
	input["__method_"] = strings.ToLower(c.Request.Method)
	input["__header_"] = c.Request.Header
	input["__is_circuit_breaker_"] = getIsCircuitBreaker(c)
	input["__jwt_"] = func() interface{} {
		return getJWTRaw(c)
	}
	input["__func_http_request_"] = c.Request
	input["__func_http_response_"] = c.Writer
	input["__binding_"] = c.ShouldBind
	input["__binding_with_"] = func(v interface{}, ct string) error {
		return c.BindWith(v, binding.Default(c.Request.Method, ct))

	}
	input["__get_request_values_"] = func() map[string]interface{} {
		c.Request.ParseForm()
		data := make(map[string]interface{})
		query := c.Request.URL.Query()
		for k, v := range query {
			switch len(v) {
			case 1:
				data[k] = v[0]
			default:
				data[k] = strings.Join(v, ",")
			}
		}
		forms := c.Request.PostForm
		for k, v := range forms {
			switch len(v) {
			case 1:
				data[k] = v[0]
			default:
				data[k] = strings.Join(v, ",")
			}
		}

		return data
	}

	input["__func_body_get_"] = func(ch string) (string, error) {
		buff, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return "", err
		}
		nbuff, err := encoding.DecodeBytes(buff, ch)
		if err != nil {
			return "", err
		}
		return string(nbuff), nil

	}
	return input
}

type MapData map[string]interface{}

//Get 获取指定键对应的数据
func (i MapData) Get(key string) (interface{}, bool) {
	r, ok := i[key]
	return r, ok
}

//InputData 输入参数
type InputData func(key string) (string, bool)

//Get 获取指定键对应的数据
func (i InputData) Get(key string) (interface{}, bool) {
	r, ok := i(key)
	return r, ok
}

//ParamData map参数数据
type ParamData map[string]string

//Get 获取指定键对应的数据
func (i ParamData) Get(key string) (interface{}, bool) {
	r, ok := i[key]
	return r, ok
}
