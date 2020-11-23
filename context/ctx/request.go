package ctx

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

type request struct {
	ctx     context.IInnerContext
	appConf app.IAPPConf
	types.XMap
	readMapErr error
	*body
	path *rpath
	*file
}

//NewRequest 构建请求的Request
//自动对请求进行解码，响应结果进行编码。
//当指定为gbk,gb2312后,请求方式为application/x-www-form-urlencoded或application/xml、application/json时内容必须编码为指定的格式，否则会解码失败
func NewRequest(c context.IInnerContext, s app.IAPPConf, meta conf.IMeta) *request {
	rpath := NewRpath(c, s, meta)
	req := &request{
		ctx:  c,
		body: NewBody(c, rpath.GetEncoding()),
		XMap: make(map[string]interface{}),
		path: rpath,
		file: NewFile(c, meta),
	}
	req.XMap, req.readMapErr = req.body.GetBodyMap()
	return req
}

//Path 获取请求路径信息
func (r *request) Path() context.IPath {
	return r.path
}

//Path 获取请求路径信息
func (r *request) Param(key string) string {
	return r.ctx.Param(key)
}

//Bind 根据输入参数绑定对象
func (r *request) Bind(obj interface{}) error {

	//检查输入类型
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Map {
		return fmt.Errorf("输入参数非struct,map %v", val.Kind())
	}

	//获取body数据
	mp, err := r.body.GetBodyMap()
	if err != nil {
		return err
	}
	if val.Kind() == reflect.Map {
		obj = mp
		return nil
	}

	//处理数据结构转换
	var xmap types.XMap = mp
	if err := xmap.ToStruct(obj); err != nil {
		return err
	}

	//验证数据格式
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		err = fmt.Errorf("输入参数有误 %v", err)
		return err
	}
	return nil
}

//Check 检查输入参数和配置参数是否为空 @todo 各类请求的测试
func (r *request) Check(field ...string) error {
	data, _ := r.body.GetBodyMap()
	for _, key := range field {
		if v, ok := data[key]; !ok || fmt.Sprint(v) == "" {
			return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("输入参数:%s值不能为空", key))
		}
	}
	return nil
}

//GetKeys 获取字段名称
func (r *request) GetKeys() []string {
	keys := make([]string, 0, 1)
	data, _ := r.body.GetBodyMap()
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

//GetMap 获取请求的参数信息
func (r *request) GetMap() (types.XMap, error) {
	return r.XMap, r.readMapErr
}

//Get 获取字段的值
func (r *request) Get(name string) (result string, ok bool) {
	result, ok = r.XMap.MustString(name)
	return result, ok

}

//GetString 从对象中获取数据值，如果不是字符串则返回空
func (r *request) GetString(name string, def ...string) string {
	value, _ := r.Get(name)
	return types.GetString(value, def...)
}

//GetInt 从对象中获取数据值，如果不是字符串则返回0
func (r *request) GetInt(name string, def ...int) int {
	value, _ := r.Get(name)
	return types.GetInt(value, def...)
}

//GetInt32 从对象中获取数据值，如果不是字符串则返回0
func (r *request) GetInt32(name string, def ...int32) int32 {
	value, _ := r.Get(name)
	return types.GetInt32(value, def...)
}

//GetInt64 从对象中获取数据值，如果不是字符串则返回0
func (r *request) GetInt64(name string, def ...int64) int64 {
	value, _ := r.Get(name)
	return types.GetInt64(value, def...)
}

//GetFloat32 从对象中获取数据值，如果不是字符串则返回0
func (r *request) GetFloat32(name string, def ...float32) float32 {
	value, _ := r.Get(name)
	return types.GetFloat32(value, def...)
}

//GetFloat64 从对象中获取数据值，如果不是字符串则返回0
func (r *request) GetFloat64(name string, def ...float64) float64 {
	value, _ := r.Get(name)
	return types.GetFloat64(value, def...)
}

//GetDecimal 获取类型为Decimal的值
func (r *request) GetDecimal(name string, def ...types.Decimal) types.Decimal {
	value, _ := r.Get(name)
	return types.GetDecimal(value, def...)
}

//GetBool 从对象中获取bool类型值，表示为true的值有：1, t, T, true, TRUE, True, YES, yes, Yes, Y, y, ON, on, On
func (r *request) GetBool(name string, def ...bool) bool {
	value, _ := r.Get(name)
	return types.GetBool(value, def...)
}

//GetDatetime 获取时间字段
func (r *request) GetDatetime(name string, format ...string) (time.Time, error) {
	value, _ := r.Get(name)
	return types.GetDatetime(value, format...)
}
func (r *request) IsEmpty(name string) bool {
	_, ok := r.Get(name)
	return ok
}

//GetPlayload 获取trace信息 //@fix GetTrace 改为的GetPlayload @hj
func (r *request) GetPlayload() string {
	raw, _ := r.GetRawBody()
	return string(raw)
}

//GetHeader 获取请求头信息
func (r *request) GetHeader(key string) string {
	return strings.Join(r.GetHeaders()[key], ",")
}

//GetHeaders 获取请求的header
func (r *request) GetHeaders() http.Header {
	return r.ctx.GetHeaders()
}

//GetHeaders 获取请求的header
func (r *request) GetCookies() map[string]string {
	out := make(map[string]string)
	cookies := r.ctx.GetCookies()
	for _, cookie := range cookies {
		out[cookie.Name] = cookie.Value
	}
	return out
}

//GetCookie 获取cookie信息
func (r *request) GetCookie(name string) string {
	cookies := r.ctx.GetCookies()
	for _, cookie := range cookies {
		if name == cookie.Name {
			return cookie.Value
		}
	}
	return ""
}
