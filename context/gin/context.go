package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

var _ context.IContext = &GinCtx{}

//GinCtx gin.context
type GinCtx struct {
	context     *gin.Context
	log         logger.ILogger
	request     *request
	response    *response
	user        *user
	server      server.IServerConf
	application context.IApplication
}

//NewGinCtx 构建基于gin.Context的上下文
func NewGinCtx(c *gin.Context) *GinCtx {
	ctx := &GinCtx{
		context:  c,
		user:     &user{Context: c},
		response: &response{Context: c},
	}
	ctx.request = newRequest(c)
	ctx.log = logger.GetSession(ctx.server.GetMainConf().GetServerName(), ctx.User().GetRequestID())
	return ctx
}

//Request 获取请求对象
func (c *GinCtx) Request() context.IRequest {
	return c.request
}

//Response 获取响应对象
func (c *GinCtx) Response() context.IResponse {
	return c.response
}

//User 获取用户相关信息
func (c *GinCtx) User() context.IUser {
	return c.user
}

//Log 获取日志组件
func (c *GinCtx) Log() logger.ILogger {
	return c.log
}

//Server 获取服务器配置
func (c *GinCtx) Server() server.IServerConf {
	return c.server
}

//Application 获取应用配置
func (c *GinCtx) Application() context.IApplication {
	return c.application
}

//Close 关闭并释放所有资源
func (c *GinCtx) Close() {

}

//Next 执行下一个中间件(中间件调用)
func (c *GinCtx) Next() {
	c.context.Next()
}
