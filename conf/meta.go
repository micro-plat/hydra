package conf

import (
	"github.com/micro-plat/lib4go/types"
)

//IMeta 无数据接口
type IMeta = types.IXMap

//Meta 元数据
type Meta = types.XMap

//NewMeta 构建元数据
func NewMeta() Meta {
	return make(map[string]interface{})
}
