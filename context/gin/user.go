package gin

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

var xRequestID = "X-Request-Id"

var _ context.IUser = &user{}

type user struct {
	*gin.Context
	requestID string
	jwtToken  interface{}
}

//GetRequestID 获取请求编号
func (c *user) GetRequestID() string {
	c.requestID = types.GetString(c.Context.GetHeader(xRequestID), c.requestID)
	c.requestID = types.GetString(c.requestID, utility.GetGUID()[0:9])
	return c.requestID
}

//GetClientIP 获取客户端IP地址
func (c *user) GetClientIP() string {
	return c.Context.ClientIP()
}

//SaveJwt  保存用户登录登录成功的jwt信息
func (c *user) SaveJwt(v interface{}) {
	c.jwtToken = v
}

//GetJwt 获取已设置的jwt信息
func (c *user) GetJwt() interface{} {
	return c.jwtToken
}

//BindJwt 绑定jwt信息
func (c *user) BindJwt(out interface{}) {
	if c.jwtToken == nil || reflect.ValueOf(c.jwtToken).IsNil() {
		panic(errs.NewError(401, "请求中未包含jwt,用户未登录"))
	}
	switch v := c.jwtToken.(type) {
	case func() interface{}:
		r := v()
		if r == nil {
			panic(errs.NewError(401, "请求中未包含jwt,用户未登录"))
		}
		buff, err := json.Marshal(r)
		if err != nil {
			panic(fmt.Errorf("将用户jwt信息转换为json失败:%w", err))
		}
		if err := json.Unmarshal(buff, &out); err != nil {
			panic(fmt.Errorf("将用户jwt信息反序化为对象时失败:%w", err))
		}
	case string:
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			panic(fmt.Errorf("将用户jwt信息反序化为对象时失败:%w", err))
		}
	default:
		buff, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Errorf("将用户jwt信息转换为json失败:%w", err))
		}
		if err := json.Unmarshal(buff, &out); err != nil {
			panic(fmt.Errorf("将用户jwt信息反序化为对象时失败:%w", err))
		}
	}
}
