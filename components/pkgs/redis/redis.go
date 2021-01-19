package redis

import (
	"time"

	"github.com/go-redis/redis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/lib4go/types"
)

//Client redis客户端
type Client struct {
	redis.UniversalClient
	opt *varredis.Redis
}

//NewByOpts 构建客户端
func NewByOpts(opts ...varredis.Option) (r *Client, err error) {
	redisOpts := varredis.New("")
	for i := range opts {
		opts[i](redisOpts)
	}
	return NewByConfig(redisOpts)
}

//NewByConfig 构建客户端
func NewByConfig(config *varredis.Redis) (r *Client, err error) {
	r = &Client{
		opt: config,
	}

	r.opt.DialTimeout = types.DecodeInt(r.opt.DialTimeout, 0, 3, r.opt.DialTimeout)
	r.opt.ReadTimeout = types.DecodeInt(r.opt.ReadTimeout, 0, 3, r.opt.ReadTimeout)
	r.opt.WriteTimeout = types.DecodeInt(r.opt.WriteTimeout, 0, 3, r.opt.WriteTimeout)
	r.opt.PoolSize = types.DecodeInt(r.opt.PoolSize, 0, 10, r.opt.PoolSize)
	ropts := &redis.UniversalOptions{
		Addrs:        r.opt.Addrs,
		Password:     r.opt.Password,
		DB:           r.opt.DbIndex,
		DialTimeout:  time.Duration(r.opt.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(r.opt.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(r.opt.WriteTimeout) * time.Second,
		PoolSize:     r.opt.PoolSize,
	}
	r.UniversalClient = redis.NewUniversalClient(ropts)
	_, err = r.UniversalClient.Ping().Result()
	return
}

//GetAddrs GetAddrs
func (c *Client) GetAddrs() []string {
	return c.opt.Addrs
}
