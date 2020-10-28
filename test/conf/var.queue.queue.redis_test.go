package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/conf/vars/queue/redis"
	"github.com/micro-plat/hydra/test/assert"
)

func TestQueueRedisNew(t *testing.T) {

	tests := []struct {
		name    string
		address string
		opts    []redis.Option
		want    *redis.Redis
	}{
		{
			name:    "测试新增-无option",
			address: "192.168.5.79:6379",
			want: &redis.Redis{
				Queue:        &queue.Queue{Proto: "redis"},
				Address:      "192.168.5.79:6379",
				DbIndex:      1,
				DialTimeout:  10,
				ReadTimeout:  10,
				WriteTimeout: 10,
				PoolSize:     10,
			},
		},
		{
			name:    "测试新增-WithDbIndex",
			address: "192.168.5.79:6379",
			opts: []redis.Option{
				redis.WithDbIndex(2),
			},
			want: &redis.Redis{
				Address:      "192.168.5.79:6379",
				Queue:        &queue.Queue{Proto: "redis"},
				DbIndex:      2,
				DialTimeout:  10,
				ReadTimeout:  10,
				WriteTimeout: 10,
				PoolSize:     10,
			},
		},
		{
			name:    "测试新增-WithTimeout",
			address: "192.168.5.79:6379",
			opts: []redis.Option{
				redis.WithDbIndex(2),
				redis.WithTimeout(11, 22, 33),
			},
			want: &redis.Redis{
				Address:      "192.168.5.79:6379",
				Queue:        &queue.Queue{Proto: "redis"},
				DbIndex:      2,
				DialTimeout:  11,
				ReadTimeout:  22,
				WriteTimeout: 33,
				PoolSize:     10,
			},
		},
		{
			name:    "测试新增-WithTimeout",
			address: "192.168.5.79:6379",
			opts: []redis.Option{
				redis.WithDbIndex(2),
				redis.WithTimeout(11, 22, 33),
				redis.WithPoolSize(40),
			},
			want: &redis.Redis{
				Address:      "192.168.5.79:6379",
				Queue:        &queue.Queue{Proto: "redis"},
				DbIndex:      2,
				DialTimeout:  11,
				ReadTimeout:  22,
				WriteTimeout: 33,
				PoolSize:     40,
			},
		},
	}
	for _, tt := range tests {
		got := redis.New(tt.address, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
