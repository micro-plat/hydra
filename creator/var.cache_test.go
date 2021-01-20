package creator

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name string
		args map[string]map[string]interface{}
		want *Varcache
	}{
		{name: "1. 初始化cache对象", args: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}},
			want: &Varcache{vars: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}}}},
	}
	for _, tt := range tests {
		got := NewCache(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarcache_Redis(t *testing.T) {

	type args struct {
		name    string
		address string
		opts    []cacheredis.Option
	}
	tests := []struct {
		name    string
		fields  *Varcache
		args    args
		want    vars
		wantErr string
	}{
		{name: "1. configname是空", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "redis", address: "address"},
			want: map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"redis": cacheredis.New("address")}}},
	}

	for _, tt := range tests {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					assert.Equal(t, tt.wantErr, types.GetString(e), tt.name+",err")
				}
			}()

			got := tt.fields.Redis(tt.args.name, tt.args.address, tt.args.opts...)
			assert.Equal(t, tt.want, got, tt.name)
		}()
	}
}

func TestVarcache_GoCache(t *testing.T) {
	type args struct {
		name string
		q    gocache.Option
	}
	tests := []struct {
		name   string
		fields *Varcache
		args   args
		want   vars
	}{
		{name: "1. 初始化GoCache对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "gocache", q: gocache.WithCleanupInterval(10)},
			want: map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"gocache": gocache.New(gocache.WithCleanupInterval(10))}}},
	}
	for _, tt := range tests {
		got := tt.fields.GoCache(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarcache_Memcache(t *testing.T) {
	type args struct {
		name string
		addr string
		opts memcached.Option
	}
	tests := []struct {
		name   string
		fields *Varcache
		args   args
		want   vars
	}{
		{name: "1. 初始化Memcache对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "memcached", opts: memcached.WithTimeout(10)},
			want: map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"memcached": memcached.New("", memcached.WithTimeout(10))}}},
	}
	for _, tt := range tests {
		got := tt.fields.Memcache(tt.args.name, tt.args.addr, tt.args.opts)
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
		repeat *args
		want   vars
	}{
		{name: "1. 初始化空Custom对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "", q: map[string]interface{}{}},
			want: map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"": map[string]interface{}{}}}},
		{name: "2. 初始化自定义Custom对象", fields: NewCache(map[string]map[string]interface{}{}), args: args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			want: map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"sss": "sdfdsfsdf"}}}},
		{name: "3. 重复初始化Custom对象", fields: NewCache(map[string]map[string]interface{}{}),
			args:   args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			repeat: &args{name: "customer", q: map[string]interface{}{"www": "54dfdf"}},
			want:   map[string]map[string]interface{}{cache.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"www": "54dfdf"}}}},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		if tt.repeat != nil {
			got = tt.fields.Custom(tt.repeat.name, tt.repeat.q)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}
