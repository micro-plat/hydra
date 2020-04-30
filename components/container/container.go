package container

import (
	"github.com/micro-plat/hydra/registry/conf"
	"github.com/qxnw/lib4go/concurrent/cmap"
)

//ICloser 关闭
type ICloser interface {
	Close() error
}

//IContainer 组件容器
type IContainer interface {
	Conf() conf.IVarConf
	GetOrCreate(name string, creator func(i ...interface{}) (interface{}, error)) (interface{}, error)
	ICloser
}

//Container 容器用于缓存公共组件
type Container struct {
	conf  conf.IVarConf
	cache cmap.ConcurrentMap
}

//NewContainer 构建容器
func NewContainer(c conf.IVarConf) *Container {
	return &Container{
		conf:  c,
		cache: cmap.New(8),
	}

}

//Conf 获取配置信息
func (c *Container) Conf() conf.IVarConf {
	return c.conf
}

//GetOrCreate 获取指定名称的组件，不存在时自动创建
func (c *Container) GetOrCreate(name string, creator func(i ...interface{}) (interface{}, error)) (interface{}, error) {
	_, obj, err := c.cache.SetIfAbsentCb(name, creator)
	return obj, err

}

//Close 释放组件资源
func (c *Container) Close() error {
	c.cache.RemoveIterCb(func(key string, v interface{}) bool {
		if closer, ok := v.(ICloser); ok {
			closer.Close()
		}
		return true
	})
	return nil
}
