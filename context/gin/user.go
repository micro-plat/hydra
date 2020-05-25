package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

var xRequestID = "X-Request-Id"

var _ context.IUser = &user{}

type user struct {
	*gin.Context
	requestID string
	auth      *auth
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

//Auth 用户认证信息
func (c *user) Auth() IAuth {
	return c.auth
}
