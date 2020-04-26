package context

import (
	"fmt"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var globals cmap.ConcurrentMap

func init() {
	globals = cmap.New(3)
}

//IUUID 全局功能接口
type IUUID interface {
	GetString(pre ...string) string
	Get() int64
}

//NewUUID 获取使用雪花算法得到的集群唯一编号
func NewUUID(container interface{}) IUUID {
	ct, ok := container.(IContainer)
	if !ok {
		panic(fmt.Errorf("输入参数container必须实现了IContainer接口"))
	}
	_, sf, err := globals.SetIfAbsentCb(ct.GetServerPubRootPath(), func(i ...interface{}) (interface{}, error) {
		return NewSnowflake(ct.GetServerPubRootPath(),
			ct.GetClusterID(),
			ct.GetRegistry())

	})
	if err != nil {
		panic(err)
	}
	return sf.(*Snowflake)
}
