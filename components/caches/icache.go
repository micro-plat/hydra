package caches

import "github.com/micro-plat/hydra/components/caches/cache"

//ICache 缓存接口
type ICache = cache.ICache

//IComponentCache Component Cache
type IComponentCache interface {
	GetRegularCache(names ...string) (c ICache)
	GetCache(names ...string) (c ICache, err error)
}
