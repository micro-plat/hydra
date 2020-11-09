package confvars

import (
	"testing"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/vars/cache"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
	cacheredis "github.com/micro-plat/hydra/conf/vars/cache/redis"
	"github.com/micro-plat/hydra/conf/vars/redis"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name string
		args map[string]map[string]interface{}
		want *Varcache
	}{
		{name: "初始化cacha对象", args: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}},
			want: &Varcache{vars: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}}}},
	}
	for _, tt := range tests {
		got := NewCache(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarcache_Redis(t *testing.T) {

	newRedis := cacheredis.New(cacheredis.WithConfigName("address"))
	newRedis.Redis = redis.New([]string{"192.196.0.1"})
	type args struct {
		name string
		q    *cacheredis.Redis
	}
	tests := []struct {
		name    string
		fields  *Varcache
		args    args
		want    *Varcache
		wantErr string
	}{
		{name: "configname是空", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "redis", q: cacheredis.New(cacheredis.WithAddrs("address"))},
			want: NewCache(map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"redis": cacheredis.New(cacheredis.WithAddrs("address"))}})},
		{name: "configname不为空,无redis节点", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "redis", q: cacheredis.New(cacheredis.WithConfigName("address"))},
			want: nil, wantErr: "请确认已配置/var/redis"},
		{name: "configname不为空,无configname节点", fields: NewCache(map[string]map[string]interface{}{"redis": map[string]interface{}{}}), args: args{name: "redis", q: cacheredis.New(cacheredis.WithConfigName("address"))},
			want: nil, wantErr: "请确认已配置/var/redis/address"},
		{name: "configname不为空,configname节点", fields: NewCache(map[string]map[string]interface{}{"redis": map[string]interface{}{"address": redis.New([]string{"192.196.0.1"})}}),
			args: args{name: "redis", q: cacheredis.New(cacheredis.WithConfigName("address"))},
			want: NewCache(map[string]map[string]interface{}{"redis": map[string]interface{}{"address": redis.New([]string{"192.196.0.1"})},
				cache.TypeNodeName: map[string]interface{}{"redis": newRedis}}), wantErr: ""},
	}

	for _, tt := range tests {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					assert.Equal(t, tt.wantErr, types.GetString(e), tt.name+",err")
				}
			}()

			got := tt.fields.Redis(tt.args.name, tt.args.q)
			assert.Equal(t, tt.want, got, tt.name)
		}()
	}
}

func TestVarcache_GoCache(t *testing.T) {
	type args struct {
		name string
		q    *gocache.GoCache
	}
	tests := []struct {
		name   string
		fields *Varcache
		args   args
		want   *Varcache
	}{
		{name: "初始化对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "gocache", q: gocache.New(gocache.WithCleanupInterval(10))},
			want: NewCache(map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"gocache": gocache.New(gocache.WithCleanupInterval(10))}})},
	}
	for _, tt := range tests {
		got := tt.fields.GoCache(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarcache_Memcache(t *testing.T) {
	type args struct {
		name string
		q    *memcached.Memcache
	}
	tests := []struct {
		name   string
		fields *Varcache
		args   args
		want   *Varcache
	}{
		{name: "初始化对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "memcached", q: memcached.New(memcached.WithTimeout(10))},
			want: NewCache(map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"memcached": memcached.New(memcached.WithTimeout(10))}})},
	}
	for _, tt := range tests {
		got := tt.fields.Memcache(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarcache_Custom(t *testing.T) {
	type args struct {
		name string
		q    interface{}
	}
	tests := []struct {
		name   string
		fields *Varcache
		args   args
		want   *Varcache
	}{
		{name: "初始化对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			want: NewCache(map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"sss": "sdfdsfsdf"}}})},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
