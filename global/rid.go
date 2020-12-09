package global

import (
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/security/md5"
)

//RID 请求唯一标识管理
var RID = &rid{cache: cmap.New(16)}

type rid struct {
	cache cmap.ConcurrentMap
}

//GetXRequestID 获取用户请求编号
func (i *rid) GetXRequestID() string {
	// id := utility.GetGUID()[0:9]
	id := md5.Encrypt(GetGoroutineID())
	if v, ok := i.cache.Get(id); ok {
		return v.(string)
	}
	return id
}

//Add 添加新编号
func (i *rid) Add(nid string) {
	i.cache.Set(md5.Encrypt(GetGoroutineID()), nid)
}

//Remove 移除当前用户编号
func (i *rid) Remove() {
	i.cache.Remove(md5.Encrypt(GetGoroutineID()))
}
