package context

import (
	"fmt"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/micro-plat/lib4go/utility"
)

type IData interface {
	//Set(string, string)
	Get(string) (string, error)
	//Each(func(string, string))
}

//Request 输入参数
type Request struct {
	Form           *inputParams
	QueryString    *inputParams
	Param          *inputParams
	Setting        *inputParams
	CircuitBreaker *circuitBreakerParam //熔断处理
	Http           *httpRequest
	Ext            *extParams
}

func newRequest() *Request {
	return &Request{
		QueryString:    &inputParams{},
		Form:           &inputParams{},
		Param:          &inputParams{},
		Setting:        &inputParams{},
		CircuitBreaker: &circuitBreakerParam{inputParams: &inputParams{}},
		Http:           &httpRequest{},
		Ext:            &extParams{},
	}
}

func (r *Request) reset(queryString IData, form IData, param IData, setting IData, ext map[string]interface{}) {
	r.QueryString.data = queryString
	r.Form.data = form
	r.Param.data = param
	r.Setting.data = setting
	r.CircuitBreaker.inputParams.data = setting
	r.CircuitBreaker.ext = ext
	r.Ext.ext = ext
	r.Http.ext = ext

}

//Bind 根据输入参数绑定对象
func (r *Request) Bind(obj interface{}) error {
	f := r.Ext.GetBindingFunc()
	if err := f(obj); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		err = fmt.Errorf("输入参数有误 %v", err)
		return err
	}
	return nil
}

//BindWith 根据输入参数绑定对象
func (r *Request) BindWith(obj interface{}, contentType string) error {
	f := r.Ext.GetBindWithFunc()
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
func (r *Request) Check(checker map[string][]string) (int, error) {
	for _, field := range checker["input"] {
		if err := r.Form.Check(field); err == nil {
			continue
		}
		if err := r.QueryString.Check(field); err != nil {
			return ERR_NOT_ACCEPTABLE, fmt.Errorf("输入参数:%v", err)
		}
	}
	if err := r.Setting.Check(checker["setting"]...); err != nil {
		return ERR_NOT_EXTENDED, fmt.Errorf("配置参数:%v", err)
	}
	return 0, nil
}

//Body2Input 根据编码格式解码body参数，并更新input参数
func (r *Request) Body2Input(encoding ...string) (map[string]string, error) {
	body, err := r.Ext.GetBody(encoding...)
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
func (r *Request) Get(name string) (result string, err error) {
	if result, err = r.Form.Get(name); err == nil {
		return result, nil
	}
	if result, err = r.QueryString.Get(name); err == nil {
		return result, nil
	}
	return "", fmt.Errorf("%s值不存在", name)
}

func (r *Request) GetString(name string, p ...string) string {
	v, err := r.Get(name)
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

//GetInt 获取int数字
func (r *Request) GetInt(name string, p ...int) int {
	value, err := r.Get(name)
	var v int
	if err == nil {
		v, err = strconv.Atoi(value)
	}
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetInt64 获取int64数字
func (r *Request) GetInt64(name string, p ...int64) int64 {
	value, err := r.Get(name)
	var v int64
	if err == nil {
		v, err = strconv.ParseInt(value, 10, 64)
	}
	if err == nil {
		return v
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

//GetFloat64 获取float64数字
func (r *Request) GetFloat64(name string, p ...float64) float64 {
	value, err := r.Get(name)
	var v float64
	if err == nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	if err == nil {
		return v
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
	value, err := r.Get(name)
	var v time.Time
	if err == nil {
		v, err = time.Parse(format, value)
	}
	if err == nil {
		return v, nil
	}
	if len(p) > 0 {
		return p[0], nil
	}
	return v, err
}

//clear 清空数据
func (r *Request) clear() {
	r.QueryString.data = nil
	r.Form.data = nil
	r.Param.data = nil
	r.Setting.data = nil
	r.CircuitBreaker.inputParams.data = nil
	r.CircuitBreaker.ext = nil
	r.Ext.ext = nil
	r.Http.ext = nil
}
