package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/test/assert"
)

func TestCacheRedisNew(t *testing.T) {

	tests := []struct {
		name    string
		address string
		opts    []cacheredis.Option
		want    *cacheredis.Redis
	}{
		{
			name: "测试新增-无option",
			opts: []cacheredis.Option{
				cacheredis.WithAddrs("192.168.5.79:6379"),
			},
			want: &cacheredis.Redis{
				Cache: &cache.Cache{Proto: "redis"},
				Redis: &varredis.Redis{
					Addrs:        []string{"192.168.5.79:6379"},
					DbIndex:      0,
					DialTimeout:  10,
					ReadTimeout:  10,
					WriteTimeout: 10,
					PoolSize:     10,
				},
			},
		},
		{
			name: "测试新增-WithDbIndex",
			opts: []cacheredis.Option{
				cacheredis.WithAddrs("192.168.5.79:6379"),
				cacheredis.WithDbIndex(2),
			},
			want: &cacheredis.Redis{
				Cache: &cache.Cache{Proto: "redis"},
				Redis: &varredis.Redis{
					Addrs:        []string{"192.168.5.79:6379"},
					DbIndex:      2,
					DialTimeout:  10,
					ReadTimeout:  10,
					WriteTimeout: 10,
					PoolSize:     10,
				},
			},
		},
		{
			name: "测试新增-WithTimeout",
			opts: []cacheredis.Option{
				cacheredis.WithAddrs("192.168.5.79:6379"),
				cacheredis.WithDbIndex(2),
				cacheredis.WithTimeout(11, 22, 33),
			},
			want: &cacheredis.Redis{
				Cache: &cache.Cache{Proto: "redis"},
				Redis: &varredis.Redis{
					Addrs:        []string{"192.168.5.79:6379"},
					DbIndex:      2,
					DialTimeout:  11,
					ReadTimeout:  22,
					WriteTimeout: 33,
					PoolSize:     10,
				},
			},
		},
		{
			name: "测试新增-WithPoolSize",
			opts: []cacheredis.Option{
				cacheredis.WithAddrs("192.168.5.79:6379"),
				cacheredis.WithDbIndex(2),
				cacheredis.WithTimeout(11, 22, 33),
				cacheredis.WithPoolSize(40),
			},
			want: &cacheredis.Redis{
				Cache: &cache.Cache{Proto: "redis"},
				Redis: &varredis.Redis{
					Addrs:        []string{"192.168.5.79:6379"},
					DbIndex:      2,
					DialTimeout:  11,
					ReadTimeout:  22,
					WriteTimeout: 33,
					PoolSize:     40,
				},
			},
		},
	}

	for _, tt := range tests {
		got := cacheredis.New(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
