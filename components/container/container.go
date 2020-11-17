package container

import (
	"fmt"
	"io"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//ICloser 关闭
type ICloser interface {
	Close() error
}

//IContainer 组件容器
type IContainer interface {
	GetOrCreate(typ string, name string, creator func(conf conf.IVarConf) (interface{}, error)) (interface{}, error)
	ICloser
}

//Container 容器用于缓存公共组件
type Container struct {
	cache cmap.ConcurrentMap
	vers  *vers
}

//NewContainer 构建容器
func NewContainer() *Container {
	c := &Container{
		cache: cmap.New(8),
		vers:  newVers(),
	}
	go c.clear()
	return c

}

//GetOrCreate 获取指定名称的组件，不存在时自动创建
func (c *Container) GetOrCreate(typ string, name string, creator func(conf conf.IVarConf) (interface{}, error)) (interface{}, error) {

	var jsconf = conf.EmptyRawConf
	vc, _ := app.Cache.GetVarConf()
	if vc != nil {
		jsconf, _ = vc.GetConf(typ, name)
	} else {
		vc = vars.EmptyVarConf
	}

	key := fmt.Sprintf("%s_%s_%d", typ, name, jsconf.GetVersion())
	_, obj, err := c.cache.SetIfAbsentCb(key, func(i ...interface{}) (interface{}, error) {
		v, err := creator(vc)
		if err != nil {
			return nil, err
		}
		c.vers.Add(typ, name, key)
		return v, nil
	})
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
func (c *Container) clear() {
	tk := time.NewTicker(time.Hour)
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-tk.C:
			c.vers.Remove(func(key string) bool {
				v, ok := c.cache.Get(key)
				if !ok {
					return true
				}
				if c, ok := v.(io.Closer); ok {
					c.Close()
				}
				return true
			})
		}
	}
}
