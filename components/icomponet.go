package components

import (
	"github.com/micro-plat/hydra/conf"
)

//IComponents 组件容器
type IComponents interface {
	GetOrCreateByConf(typeName string, nodeName string, creator func(c conf.IConf) (interface{}, error)) (interface{}, error)
	GetOrCreate(name string, creator func(i ...interface{}) (interface{}, error)) (interface{}, error)
	Get(cacheKey string) (interface{}, error)
	GetRequestID() string
}
