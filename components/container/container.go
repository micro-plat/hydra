package container

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//ICloser 关闭
type ICloser interface {
	Close() error
}

//IContainer 组件容器
type IContainer interface {
	GetOrCreate(typ string, name string, creator func(conf *conf.RawConf, keys ...string) (interface{}, error), keys ...string) (interface{}, error)
	ICloser
}

//Container 容器用于缓存公共组件
type Container struct {
	cache     cmap.ConcurrentMap
	histories *histories
}

//NewContainer 构建容器,用于管理公共组件
func NewContainer() *Container {
	c := &Container{
		cache:     cmap.New(8),
		histories: newHistories(),
	}
	go c.clear()
	return c

}

//GetOrCreate 获取指定名称的组件，不存在时自动创建
func (c *Container) GetOrCreate(typ string, name string, creator func(conf *conf.RawConf, keys ...string) (interface{}, error), keys ...string) (interface{}, error) {

	//1. 获取配置信息
	varConf, err := app.Cache.GetVarConf()
	if err != nil {
		return nil, fmt.Errorf("无法获取var.conf:%w", err)
	}
	jconf := conf.EmptyRawConf
	if varConf.Has(typ, name) {
		jconf, err = varConf.GetConf(typ, name)
		if err != nil {
			return nil, err
		}
	}

	//2. 根据配置创建组件
	key := fmt.Sprintf("%s_%s_%s_%d", typ, name, strings.Join(keys, "_"), jconf.GetVersion())
	_, obj, err := c.cache.SetIfAbsentCb(key, func(i ...interface{}) (interface{}, error) {
		nkeys := []string{}
		if len(i) > 1 {
			nkeys = i[1].([]string)
		}
		v, err := creator(i[0].(*conf.RawConf), nkeys...)
		if err != nil {
			return nil, err
		}
		c.histories.Add(fmt.Sprintf("%s_%s_%s", typ, name, strings.Join(keys, "_")), key)
		return v, nil
	}, jconf, keys)
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
			c.histories.Remove(func(key string) bool {
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
