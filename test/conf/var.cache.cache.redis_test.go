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
		opts    []redis.Option
		want    *redis.Redis
	}{
		{
			name: "测试新增-无option",
			opts: []redis.Option{
				redis.WithAddrs("192.168.5.79:6379"),
			},
			want: &redis.Redis{
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
			opts: []redis.Option{
				redis.WithAddrs("192.168.5.79:6379"),
				redis.WithDbIndex(2),
			},
			want: &redis.Redis{
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
			opts: []redis.Option{
				redis.WithAddrs("192.168.5.79:6379"),
				redis.WithDbIndex(2),
				redis.WithTimeout(11, 22, 33),
			},
			want: &redis.Redis{
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
			opts: []redis.Option{
				redis.WithAddrs("192.168.5.79:6379"),
				redis.WithDbIndex(2),
				redis.WithTimeout(11, 22, 33),
				redis.WithPoolSize(40),
			},
			want: &redis.Redis{
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
		got := redis.New(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
