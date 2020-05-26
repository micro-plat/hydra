package ctx

import (
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

var xRequestID = "X-Request-Id"

var _ context.IUser = &user{}

type user struct {
	ctx       context.IInnerContext
	requestID string
	auth      *auth
	jwtToken  interface{}
}

//GetRequestID 获取请求编号
func (c *user) GetRequestID() string {
	ids := c.ctx.GetHeaders()[xRequestID]
	c.requestID = types.GetStringByIndex(ids, 0, c.requestID)
	c.requestID = types.GetString(c.requestID, utility.GetGUID()[0:9])
	return c.requestID
}

//GetClientIP 获取客户端IP地址
func (c *user) GetClientIP() string {
	return c.ctx.ClientIP()
}

//Auth 用户认证信息
func (c *user) Auth() context.IAuth {
	return c.auth
}
