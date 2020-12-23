package ctx

import (
	r "context"
	"fmt"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx/internal"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"
	"sync"
	"time"
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
	tracer     *tracer
	cancelFunc func()
}

//NewCtx 构建基于gin.Context的上下文
func NewCtx(c context.IInnerContext, tp string) *Ctx {
	var err error
	ctx := contextPool.Get().(*Ctx)
	ctx.meta = conf.NewMeta()
	ctx.context = c
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
	ctx.tracer = newTracer(c.GetURL().Path, ctx.log, ctx.appConf)
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

//Tracer 链路跟踪器
func (c *Ctx) Tracer() context.ITracer {
	return c.tracer
}

//Invoke 调用本地服务
func (c *Ctx) Invoke(service string) interface{} {
	proto, addr, err := global.ParseProto(service)
	if err != nil {
		return err
	}
	switch proto {
	case global.ProtoRPC:
		return internal.CallRPC(c, addr)
	case global.ProtoInvoker:
		return services.Def.Call(c, addr)
	}
	return fmt.Errorf("不支持%s的服务调用%s", proto, service)

}

//Close 关闭并释放所有资源
func (c *Ctx) Close() {
	context.Del(c.user.gid) //从当前请求上下文中删除
	c.appConf = nil
	c.cancelFunc()
	c.cancelFunc = nil
	c.context = nil
	c.ctx = nil
	c.log = nil
	c.meta = nil
	c.request = nil
	c.response = nil
	c.user = nil
	c.tracer = nil

	contextPool.Put(c)
}
