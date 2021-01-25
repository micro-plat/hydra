package ctx

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

type Auth struct {
	request  interface{}
	response interface{}
	c        context.IInnerContext
}

func newAuth(c context.IInnerContext) *Auth {
	return &Auth{
		c: c,
	}
}

//Response  用户响应的认证信息
func (c *Auth) Response(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.response = v[0]
	}
	if c.response == nil {
		return c.request
	}
	return c.response
}

//Request  用户请求的认证信息
func (c *Auth) Request(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.request = v[0]
	}
	return c.request
}

//Clear  清除认证信息
func (c *Auth) Clear() {
	c.c.ClearAuth(true)
}

//Bind 绑定用户信息
func (c *Auth) Bind(out interface{}) error {

	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("输入参数非指针 %v", val.Kind())
	}

	if c.request == nil || c.isNil(c.request) {
		return errs.NewError(401, "请求中未包含用户信息,用户未登录")
	}

	switch v := c.request.(type) {
	case func() interface{}:
		r := v()
		if r == nil {
			return errs.NewError(401, "请求中未包含用户信息,用户未登录")
		}
		buff, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("将用户信息转换为json失败:%w", err)
		}
		if err := json.Unmarshal(buff, out); err != nil {
			return fmt.Errorf("将用户信息反序化为对象时失败(func):%w ;%s", err, string(buff))
		}
	case string:
		if err := json.Unmarshal([]byte(v), out); err != nil {
			return fmt.Errorf("将用户信息反序化为对象时失败(string):%w;%s", err, v)
		}
	default:
		buff, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("将用户信息转换为json失败:%w", err)
		}
		if err := json.Unmarshal(buff, out); err != nil {
			return fmt.Errorf("将用户信息反序化为对象时失败(default):%w ;%s", err, string(buff))
		}
	}
	return nil
}

func (c *Auth) isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.String {
		return i == ""
	}
	if vi.Kind() == reflect.Ptr || vi.Kind() == reflect.Func || vi.Kind() == reflect.Map || vi.Kind() == reflect.Slice {
		return vi.IsNil()
	}
	return false
}
