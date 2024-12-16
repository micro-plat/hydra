package caches

import (
	"fmt"

	"github.com/micro-plat/hydra/components/caches/cache"
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/conf"
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
	c container.IContainer
}

//NewStandardCache 创建cache
func NewStandardCache(c container.IContainer) *StandardCache {
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
	obj, err := s.c.GetOrCreate(cacheTypeNode, name, func(conf *conf.RawConf, keys ...string) (interface{}, error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", cacheTypeNode, name)
		}
		orgCache, err := cache.New(conf.GetString("proto"), string(conf.GetRaw()))
		return orgCache, err
	})
	if err != nil {
		return nil, err
	}
	return obj.(ICache), nil
}
