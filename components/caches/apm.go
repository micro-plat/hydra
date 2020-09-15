package caches

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/components/pkgs/cache"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/apm"
)

var _ ICache = &APMCache{}

type APMCache struct {
	orgCache ICache
	name     string
	proto    string
	servers  []string
}

type DBCallback func() *CallResult
type CallResult struct {
	Exists   bool
	Val      string
	Error    error
	EffCount int64
	Vals     []string
}

func NewAPMCache(name string, orgCache ICache) (obj ICache, err error) {

	//fmt.Println("NewAPMCache.1")
	cacheExt, ok := orgCache.(cache.ICacheExt)
	if !ok {
		return orgCache, nil
	}
	//fmt.Println("NewAPMCache.2")

	return &APMCache{
		name:     name,
		orgCache: orgCache,
		proto:    cacheExt.GetProto(),
		servers:  cacheExt.GetServers(),
	}, err
}

func (d *APMCache) Get(key string) (string, error) {
	callback := func() *CallResult {
		val, err := d.orgCache.Get(key)
		return &CallResult{
			Val:   val,
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Get", key, d.servers, callback)
	return result.Val, result.Error
}

func (d *APMCache) Decrement(key string, delta int64) (n int64, err error) {
	callback := func() *CallResult {
		effCount, err := d.orgCache.Decrement(key, delta)
		return &CallResult{
			EffCount: effCount,
			Error:    err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Decrement", key, d.servers, callback)
	return result.EffCount, result.Error
}

func (d *APMCache) Increment(key string, delta int64) (n int64, err error) {
	callback := func() *CallResult {
		effCount, err := d.orgCache.Increment(key, delta)
		return &CallResult{
			EffCount: effCount,
			Error:    err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Increment", key, d.servers, callback)
	return result.EffCount, result.Error
}

func (d *APMCache) Gets(key ...string) (r []string, err error) {
	callback := func() *CallResult {
		vals, err := d.orgCache.Gets(key...)
		return &CallResult{
			Vals:  vals,
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Gets", strings.Join(key, ";"), d.servers, callback)
	return result.Vals, result.Error
}
func (d *APMCache) Add(key string, value string, expiresAt int) error {
	callback := func() *CallResult {
		err := d.orgCache.Add(key, value, expiresAt)
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Add", key, d.servers, callback)
	return result.Error
}
func (d *APMCache) Set(key string, value string, expiresAt int) error {
	callback := func() *CallResult {
		err := d.orgCache.Set(key, value, expiresAt)
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Set", key, d.servers, callback)
	return result.Error
}
func (d *APMCache) Delete(key string) error {
	callback := func() *CallResult {
		err := d.orgCache.Delete(key)
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Delete", key, d.servers, callback)
	return result.Error
}
func (d *APMCache) Exists(key string) bool {
	callback := func() *CallResult {
		exists := d.orgCache.Exists(key)
		return &CallResult{
			Exists: exists,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Exists", key, d.servers, callback)
	return result.Exists
}
func (d *APMCache) Delay(key string, expiresAt int) error {
	callback := func() *CallResult {
		err := d.orgCache.Delay(key, expiresAt)
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.proto, d.name, "Cache.Delay", key, d.servers, callback)
	return result.Error
}

func (d *APMCache) Close() error {
	return d.orgCache.Close()
}

func (d *APMCache) GetProvider() string {
	return d.proto
}

func apmExecute(provider, name, operationName, url string, servers []string, callback DBCallback) *CallResult {
	ctx := context.Current()
	apmCfg := ctx.ServerConf().GetAPMConf()
	if apmCfg.Disable {
		return callback()
	}
	if !apmCfg.GetCacheEnable(name) {
		return callback()
	}
	apmCtx := ctx.APMContext()
	if apmCtx == nil {
		return callback()
	}
	rootCtx := apmCtx.GetRootCtx()
	tracer := apmCtx.GetTracer()

	peer := strings.Join(servers, ",")
	//fmt.Println("apmExecute.2-1", peer)

	span, err := tracer.CreateExitSpan(rootCtx, operationName, peer, func(header string) error {
		return nil
	})
	if err != nil {
		err = fmt.Errorf("tracer.CreateExitSpan:%+v", err)
		return callback()
	}
	//fmt.Println("apmExecute.3")
	defer span.End()
	//执行db 请求
	res := callback()
	span.SetComponent(apm.ComponentIDGOCacheClient)
	span.Tag("CacheProto", fmt.Sprintf("%s[%s]", provider, name))
	span.Tag("CacheKey", url)
	span.SetSpanLayer(apm.SpanLayer_Cache)

	return res
}
