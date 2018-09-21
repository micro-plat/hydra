package context

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/lib4go/utility"
)

type IData interface {
	Get(string) (interface{}, bool)
}

//Request 输入参数
type Request struct {
	Form           *inputParams
	QueryString    *inputParams
	Param          *inputParams
	Setting        *inputParams
	CircuitBreaker *circuitBreakerParam //熔断处理
	Http           *httpRequest
	*extParams
}

//clear 清空数据
func (r *Request) clear() {
	r.Form.Clear()
	r.QueryString.Clear()
	r.Param.Clear()
	r.Setting.Clear()
	r.CircuitBreaker.Clear()
	r.Http.Clear()
	r.extParams.Clear()
}

func newRequest() *Request {
	return &Request{
		QueryString:    &inputParams{},
		Form:           &inputParams{},
		Param:          &inputParams{},
		Setting:        &inputParams{},
		CircuitBreaker: &circuitBreakerParam{inputParams: &inputParams{}},
		Http:           &httpRequest{},
		extParams:      &extParams{},
	}
}

func (r *Request) reset(ctx *Context, queryString IData, form IData, param IData, setting IData, ext map[string]interface{}) {
	r.QueryString.data = queryString
	r.Form.data = form
	r.Param.data = param
	r.Setting.data = setting
	r.CircuitBreaker.inputParams.data = setting
	r.CircuitBreaker.ext = ext
	r.extParams.ext = ext
	r.extParams.ctx = ctx
	r.Http.ext = ext

}

//Bind 根据输入参数绑定对象
func (r *Request) Bind(obj interface{}) error {
	f := r.GetBindingFunc()
	if err := f(obj); err != nil {
		return err
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		err = fmt.Errorf("输入参数有误 %v", err)
		return err
	}
	return nil
}

//BindWith 根据输入参数绑定对象
func (r *Request) BindWith(obj interface{}, contentType string) error {
	f := r.GetBindWithFunc()
	if err := f(obj, contentType); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		err = fmt.Errorf("输入参数有误 %v", err)
		return err
	}
	return nil
}

//Check 检查输入参数和配置参数是否为空
func (r *Request) Check(field ...string) error {
	data, err := r.GetBodyMap()
	for _, fd := range field {
		if err := r.Form.Check(fd); err == nil {
			continue
		}
		if err := r.QueryString.Check(fd); err == nil {
			continue
		}
		if v, ok := data[fd]; !ok && fmt.Sprint(v) != "" {
			return fmt.Errorf("输入参数:%s值不能为空 %v", fd, err)
		}
	}
	return nil
}

//Body2Input 根据编码格式解码body参数，并更新input参数
func (r *Request) Body2Input(encoding ...string) (map[string]string, error) {
	body, err := r.GetBody(encoding...)
	if err != nil {
		return nil, err
	}
	qString, err := utility.GetMapWithQuery(body)
	if err != nil {
		return nil, err
	}
	//for k, v := range qString {
	//w.Form.data.Set(k, v)
	//}
	return qString, nil
}

//Get 获取请求数据值
func (r *Request) Get(name string) (result string, b bool) {
	if result, b = r.Form.Get(name); b {
		return result, b
	}
	if result, b = r.QueryString.Get(name); b {
		return result, b
	}
	m, err := r.GetBodyMap()
	if err != nil {
		return "", false
	}
	v, b := m[name]
	return fmt.Sprint(v), b
}

func (r *Request) GetString(name string, p ...string) string {
	v, b := r.Get(name)
	if b {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

//GetInt 获取int数字
func (r *Request) GetInt(name string, p ...int) int {
	value, b := r.Get(name)
	var v int
	var err error
	if b {
		v, err = strconv.Atoi(value)
		if err == nil {
			return v
		}
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetInt64 获取int64数字
func (r *Request) GetInt64(name string, p ...int64) int64 {
	value, b := r.Get(name)
	var v int64
	var err error
	if b {
		if strings.Contains(strings.ToUpper(value), "E+") {
			var n float64
			_, err := fmt.Sscanf(value, "%e", &n)
			if err == nil {
				return int64(n)
			}
		}
		v, err = strconv.ParseInt(value, 10, 64)
		if err == nil {
			return v
		}
	}

	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//Translate 根据输入参数[Param,Form,QueryString,Setting]
func (r *Request) Translate(format string, a bool) string {
	str, i := r.Param.Translate(format, false)
	if i == 0 {
		return str
	}

	str, i = r.Form.Translate(str, false)
	if i == 0 {
		return str
	}
	str, i = r.QueryString.Translate(str, false)
	if i == 0 {
		return str
	}

	str, _ = r.Setting.Translate(str, a)
	return str
}

//GetFloat32 获取float64数字
func (r *Request) GetFloat32(name string, p ...float32) float32 {
	value, b := r.Get(name)
	if b {
		x, err := strconv.ParseFloat(value, 32)
		if err == nil {
			return float32(x)
		}
	}

	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetFloat64 获取float64数字
func (r *Request) GetFloat64(name string, p ...float64) float64 {
	value, b := r.Get(name)
	var v float64
	var err error
	if b {
		if strings.Contains(strings.ToUpper(value), "E+") {
			var n float64
			_, err := fmt.Sscanf(value, "%e", &n)
			if err == nil {
				return n
			}
		}

		v, err = strconv.ParseFloat(value, 64)
		if err == nil {
			return v
		}
	}

	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetDataTime 获取日期时间
func (r *Request) GetDataTime(name string, p ...time.Time) (time.Time, error) {
	return r.GetDataTimeByFormat(name, "20060102150405", p...)
}

//GetDataTimeByFormat 获取日期时间
func (r *Request) GetDataTimeByFormat(name string, format string, p ...time.Time) (time.Time, error) {
	value, b := r.Get(name)
	var v time.Time
	var err error
	if b {
		v, err = time.Parse(format, value)
		if err == nil {
			return v, nil
		}
	}

	if len(p) > 0 {
		return p[0], nil
	}
	return v, err
}
