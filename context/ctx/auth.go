package ctx

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/micro-plat/lib4go/errs"
)

type auth struct {
	request  interface{}
	response interface{}
}

//Response  用户响应的认证信息
func (c *auth) Response(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.response = v[0]
	}
	return c.response
}

//Request  用户请求的认证信息
func (c *auth) Request(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.request = v[0]
	}
	return c.request
}

//Bind 绑定用户信息
func (c *auth) Bind(out interface{}) {
	if c.request == nil || reflect.ValueOf(c.request).IsNil() {
		panic(errs.NewError(401, "请求中未包含用户信息,用户未登录"))
	}
	switch v := c.request.(type) {
	case func() interface{}:
		r := v()
		if r == nil {
			panic(errs.NewError(401, "请求中未包含用户信息,用户未登录"))
		}
		buff, err := json.Marshal(r)
		if err != nil {
			panic(fmt.Errorf("将用户信息转换为json失败:%w", err))
		}
		if err := json.Unmarshal(buff, &out); err != nil {
			panic(fmt.Errorf("将用户信息反序化为对象时失败:%w", err))
		}
	case string:
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			panic(fmt.Errorf("将用户信息反序化为对象时失败:%w", err))
		}
	default:
		buff, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Errorf("将用户信息转换为json失败:%w", err))
		}
		if err := json.Unmarshal(buff, &out); err != nil {
			panic(fmt.Errorf("将用户信息反序化为对象时失败:%w", err))
		}
	}
}
