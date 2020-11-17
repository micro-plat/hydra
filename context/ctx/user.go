package ctx

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/types"
	"github.com/micro-plat/lib4go/utility"
)

var xRequestID = "X-Request-Id"

var _ context.IUser = &user{}

type user struct {
	gid       string
	meta      conf.IMeta
	ctx       context.IInnerContext
	requestID string
	auth      *Auth
	jwtToken  interface{}
}

func NewUser(ctx context.IInnerContext, gid string, meta conf.IMeta) *user {
	return &user{
		ctx:  ctx,
		gid:  gid,
		auth: &Auth{},
		meta: meta,
	}
}

//GetRequestID 获取请求编号
func (c *user) GetRequestID() string {
	ids := c.ctx.GetHeaders()[xRequestID]
	c.requestID = types.GetStringByIndex(ids, 0, c.requestID)
	if c.requestID == "" {
		c.requestID = utility.GetGUID()[0:9]
	}
	return c.requestID
}

//GetGID 获取当前处理的goroutine id
func (c *user) GetGID() string {
	return c.gid
}

//GetClientIP 获取客户端IP地址
func (c *user) GetClientIP() string {
	ip := c.ctx.ClientIP()
	if ip == "::1" || ip == "127.0.0.1" {
		return global.LocalIP()
	}
	return ip
}

//Auth 用户认证信息
func (c *user) Auth() context.IAuth {
	return c.auth
}
