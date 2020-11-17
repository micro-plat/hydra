package ctx

import (
	r "context"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/logger"
)

var _ context.IContext = &Ctx{}
var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Ctx{}
		},
	}
}

//Ctx context
type Ctx struct {
	context    context.IInnerContext
	ctx        r.Context
	meta       conf.IMeta
	log        logger.ILogger
	request    *request
	response   *response
	user       *user
	appConf    app.IAPPConf
	cancelFunc func()
	funs       *funcs
	tid        string
}

//NewCtx 构建基于gin.Context的上下文
func NewCtx(c context.IInnerContext, tp string) *Ctx {
	ctx := contextPool.Get().(*Ctx)
	ctx.meta = conf.NewMeta()
	ctx.context = c
	var err error
	ctx.appConf, err = app.Cache.GetAPPConf(tp)
	if err != nil {
		panic(err)
	}
	ctx.user = NewUser(c, context.Cache(ctx), ctx.meta)
	ctx.request = NewRequest(c, ctx.appConf, ctx.meta)
	ctx.log = logger.GetSession(ctx.appConf.GetServerConf().GetServerName(), ctx.User().GetRequestID())
	ctx.response = NewResponse(c, ctx.appConf, ctx.log, ctx.meta)
	timeout := time.Duration(ctx.appConf.GetServerConf().GetMainConf().GetInt("", 30))
	ctx.ctx, ctx.cancelFunc = r.WithTimeout(r.WithValue(r.Background(), "X-Request-Id", ctx.user.GetRequestID()), time.Second*timeout)
	ctx.funs = newFunc(ctx)
	return ctx
}

//Meta 获取元数据配置
func (c *Ctx) Meta() conf.IMeta {
	return c.meta
}

//Request 获取请求对象
func (c *Ctx) Request() context.IRequest {
	return c.request
}

//TmplFuncs 提供用于模板转换的函数表达式
func (c *Ctx) TmplFuncs() context.TFuncs {
	return c.funs.TmplFuncs()
}

// //LuaModules 提供用于模板转换的函数表达式
// func (c *Ctx) LuaModules() lua.Modules {
// 	return c.funs.LuaFuncs()
// }

//Response 获取响应对象
func (c *Ctx) Response() context.IResponse {
	return c.response
}

//Context 处理程序退出，超时等
func (c *Ctx) Context() r.Context {
	return c.ctx
}

//User 获取用户相关信息
func (c *Ctx) User() context.IUser {
	return c.user
}

//Log 获取日志组件
func (c *Ctx) Log() logger.ILogger {
	return c.log
}

//APPConf 获取服务器配置
func (c *Ctx) APPConf() app.IAPPConf {
	return c.appConf
}

//Close 关闭并释放所有资源
func (c *Ctx) Close() {
	context.Del(c.tid) //从当前请求上下文中删除
	c.appConf = nil
	c.cancelFunc()
	c.cancelFunc = nil
	c.context = nil
	c.ctx = nil
	c.funs = nil
	c.log = nil
	c.meta = nil
	c.request = nil
	c.response = nil
	c.tid = ""
	c.user = nil

	contextPool.Put(c)
}
