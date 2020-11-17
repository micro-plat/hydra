package context

import (
	"context"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

// var ctxMap sync.Map

var ctxMap = cmap.New(6)

//Cache 将当前上下文配置保存到当前线程编号对应的缓存
func Cache(s IContext) string {
	tid := global.GetGoroutineID()
	ctxMap.SetIfAbsent(tid, s)
	return tid
}

//GetContextWithDefault 获取可用的context.Context
func GetContextWithDefault() context.Context {
	if c, ok := ctxMap.Get(global.GetGoroutineID()); ok {
		return c.(IContext).Context()
	}
	return context.WithValue(context.Background(), "X-Request-Id", global.Def.Log().GetSessionID())
}

//Current 从缓存中获取请求上下文配置
func Current(g ...string) IContext {
	gid := types.GetStringByIndex(g, 0)
	if gid == "" {
		gid = global.GetGoroutineID()
	}
	if c, ok := ctxMap.Get(gid); ok {
		return c.(IContext)
	}
	panic("未获取到当前线程的请求上下文")
}

//Del 删除当前线程的请求上下文缓存
func Del(tid string) {
	ctxMap.Remove(tid)
}
