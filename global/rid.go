package global

import (
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/utility"
)

//RID 请求唯一标识管理
var RID = &rid{cache: cmap.New(16)}

type rid struct {
	cache cmap.ConcurrentMap
}

//GetXRequestID 获取用户请求编号
func (i *rid) GetXRequestID() string {
	_, v, _ := i.cache.SetIfAbsentCb(GetGoroutineID(), func(v ...interface{}) (interface{}, error) {
		return utility.GetGUID()[:9], nil
	})
	return v.(string)
}

//Remove 移除当前用户编号
func (i *rid) Remove() {
	i.cache.Remove(GetGoroutineID())
}
