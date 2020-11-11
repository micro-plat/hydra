package context

import (
	"bytes"
	"context"
	"runtime"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func getGID() string {
	b := make([]byte, 64)
	b = b[:32]
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

// var ctxMap sync.Map

var ctxMap = cmap.New(6)

//Cache 将当前上下文配置保存到当前线程编号对应的缓存
func Cache(s IContext) string {
	tid := getGID()
	ctxMap.SetIfAbsent(tid, s)
	return tid
}

//GetContextWithDefault 获取可用的context.Context
func GetContextWithDefault() context.Context {
	if c, ok := ctxMap.Get(getGID()); ok {
		return c.(IContext).Context()
	}
	return context.WithValue(context.Background(), "X-Request-Id", global.Def.Log().GetSessionID())
}

//Current 从缓存中获取请求上下文配置
func Current() IContext {
	if c, ok := ctxMap.Get(getGID()); ok {
		return c.(IContext)
	}
	panic("未获取到当前线程的请求上下文")
}

//Del 删除当前线程的请求上下文缓存
func Del(tid string) {
	ctxMap.Remove(tid)
}
