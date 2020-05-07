package gin

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

var _ context.IContext = &GinCtx{}
var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &GinCtx{}
		},
	}
}

//GinCtx gin.context
type GinCtx struct {
	context  *gin.Context
	log      logger.ILogger
	request  *request
	response *response
	user     *user
	server   server.IServerConf
	tid      uint64
}

//NewGinCtx 构建基于gin.Context的上下文
func NewGinCtx(c *gin.Context, tp string) *GinCtx {
	ctx := contextPool.Get().(*GinCtx)
	ctx.context = c
	ctx.server = application.Current().Server(tp)
	ctx.user = &user{Context: c}
	ctx.response = &response{Context: c}
	ctx.request = newRequest(c)
	ctx.log = logger.GetSession(ctx.server.GetMainConf().GetServerName(), ctx.User().GetRequestID())
	ctx.tid = context.Cache(ctx) //保存到缓存中
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

//Component 获取组件
func (c *GinCtx) Component() components.IComponent {
	return components.Def
}

//Close 关闭并释放所有资源
func (c *GinCtx) Close() {
	context.Del(c.tid) //从当前请求上下文中删除
	contextPool.Put(c)
}
