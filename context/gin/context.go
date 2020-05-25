package gin

import (
	r "context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
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
	context    *gin.Context
	ctx        r.Context
	log        logger.ILogger
	request    *request
	response   *response
	user       *user
	serverConf server.IServerConf
	tid        uint64
}

//NewGinCtx 构建基于gin.Context的上下文
func NewGinCtx(c *gin.Context, tp string) *GinCtx {
	ctx := contextPool.Get().(*GinCtx)
	ctx.context = c
	ctx.serverConf = global.Current().Server(tp)
	ctx.user = &user{Context: c}
	ctx.response = &response{Context: c, conf: ctx.serverConf}
	ctx.request = newRequest(c)
	ctx.log = logger.GetSession(ctx.serverConf.GetMainConf().GetServerName(), ctx.User().GetRequestID())
	ctx.tid = context.Cache(ctx) //保存到缓存中
	timeout := time.Duration(ctx.serverConf.GetMainConf().GetMainConf().GetInt("", 30))
	ctx.ctx, _ = r.WithTimeout(r.WithValue(r.Background(), "X-Request-Id", ctx.user.GetRequestID()), time.Second*timeout)
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

//Context 处理程序退出，超时等
func (c *GinCtx) Context() r.Context {
	return c.ctx
}

//User 获取用户相关信息
func (c *GinCtx) User() context.IUser {
	return c.user
}

//Log 获取日志组件
func (c *GinCtx) Log() logger.ILogger {
	return c.log
}

//ServerConf 获取服务器配置
func (c *GinCtx) ServerConf() server.IServerConf {
	return c.serverConf
}

//Close 关闭并释放所有资源
func (c *GinCtx) Close() {
	context.Del(c.tid) //从当前请求上下文中删除
	contextPool.Put(c)
}
