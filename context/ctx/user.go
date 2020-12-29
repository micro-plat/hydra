package ctx

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
)

var _ context.IUser = &user{}

type user struct {
	gid string
	conf.IMeta
	ctx       context.IInnerContext
	requestID string
	auth      *Auth
	jwtToken  interface{}
}

//NewUser 用户信息
func NewUser(ctx context.IInnerContext, meta conf.IMeta) *user {
	u := &user{
		ctx:   ctx,
		auth:  &Auth{},
		IMeta: meta,
	}
	u.gid = u.GetRequestID()
	return u
}

//GetRequestID 获取请求编号
func (c *user) GetRequestID() string {
	if ids, ok := c.ctx.GetHeaders()[context.XRequestID]; ok {
		global.RID.Add(ids[0])
	}
	return global.RID.GetXRequestID()
}

//GetGID 获取当前处理的goroutine id
func (c *user) GetGID() string {
	return c.gid
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
