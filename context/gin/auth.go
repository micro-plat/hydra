package gin

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

//Save  保存用户的登录信息
func (c *auth) Save(v interface{}) {
	c.response = v
}

//Cache  将用户登录信息缓存到组件
func (c *auth) Cache(v interface{}) {
	c.request = v
}

//Get 获取请求的用户信息
func (c *auth) Get() interface{} {
	return c.request
}

//Bind 绑定用户信息
func (c *auth) Bind(out interface{}) {
	if c.auth == nil || reflect.ValueOf(c.auth).IsNil() {
		panic(errs.NewError(401, "请求中未包含用户信息,用户未登录"))
	}
	switch v := c.auth.(type) {
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
