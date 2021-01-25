package ctx

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/utility"
)

var _ context.IUser = &user{}

type user struct {
	conf.IMeta
	ctx      context.IInnerContext
	traceID  string
	auth     *Auth
	jwtToken interface{}
}

//NewUser 用户信息
func NewUser(ctx context.IInnerContext, meta conf.IMeta) *user {
	u := &user{
		ctx:   ctx,
		auth:  newAuth(ctx),
		IMeta: meta,
	}
	if ids, ok := ctx.GetHeaders()[context.XRequestID]; ok && len(ids) > 0 && ids[0] != "" {
		u.traceID = ids[0]
	} else {
		u.traceID = utility.GetGUID()[:9]
	}
	return u
}

//GetTraceID 获取链路跟踪编号
func (c *user) GetTraceID() string {
	return c.traceID
}

//GetUserName 获取用户名(basic认证启动后有效)
func (c *user) GetUserName() string {
	return c.GetString(context.UserName)
}

//GetClientIP 获取客户端IP地址
func (c *user) GetClientIP() string {
	ip := c.ctx.ClientIP()
	if ip == "" || ip == "::1" || ip == "127.0.0.1" {
		return global.LocalIP()
	}
	return ip
}

//Auth 用户认证信息
func (c *user) Auth() context.IAuth {
	return c.auth
}
