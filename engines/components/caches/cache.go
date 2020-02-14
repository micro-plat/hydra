package caches

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/cache"
	"github.com/micro-plat/lib4go/types"
)

const (
	//cacheTypeNode 缓存在var配置中的类型名称
	cacheTypeNode = "cache"

	//cacheNameNode 缓存名称在var配置中的末节点名称
	cacheNameNode = "cache"
)

//StandardCache cache
type StandardCache struct {
	c components.IComponentContainer
}

//NewStandardCache 创建cache
func NewStandardCache(c components.IComponentContainer) *StandardCache {
	return &StandardCache{c: c}
}

//GetRegularCache 获取正式的没有异常缓存实例
func (s *StandardCache) GetRegularCache(names ...string) (c ICache) {
	c, err := s.GetCache(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetCache 获取缓存操作对象
func (s *StandardCache) GetCache(names ...string) (c ICache, err error) {
	name := types.GetStringByIndex(names, 0, cacheNameNode)
	obj, err := s.c.GetOrCreateByConf(cacheTypeNode, name, func(c conf.IConf) (interface{}, error) {
		return cache.NewCache(c.GetString("proto"), string(c.GetRaw()))
	})
	if err != nil {
		return nil, err
	}
	return obj.(ICache), nil
}
