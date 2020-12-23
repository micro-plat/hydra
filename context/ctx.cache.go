package context

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

// var ctxMap sync.Map
var ctxMap = cmap.New(6)

//Cache 将当前上下文配置保存到当前线程编号对应的缓存
func Cache(s IContext) string {
	tid := global.RID.GetXRequestID()
	ctxMap.SetIfAbsent(tid, s)
	return tid
}

//Current 从缓存中获取请求上下文配置
func Current(g ...string) IContext {
	if v, ok := GetContext(g...); ok {
		return v
	}
	panic("未获取到当前线程的请求上下文")
}

//GetContext 获取当前context
func GetContext(g ...string) (IContext, bool) {
	gid := types.GetStringByIndex(g, 0)
	if gid == "" {
		gid = global.RID.GetXRequestID()
	}
	if c, ok := ctxMap.Get(gid); ok {
		return c.(IContext), true
	}
	return nil, false
}

//Del 删除当前线程的请求上下文缓存
func Del(id string) {
	ctxMap.Remove(id)
	global.RID.Remove()
}
