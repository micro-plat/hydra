package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/micro-plat/lib4go/types"
)

//Client redis客户端
type Client struct {
	redis.UniversalClient
	opt *option
}

//New 构建客户端
func New(opts ...Option) (r *Client, err error) {
	r = &Client{
		opt: &option{},
	}
	for _, opt := range opts {
		opt(r.opt)
	}

	r.opt.DialTimeout = types.DecodeInt(r.opt.DialTimeout, 0, 3, r.opt.DialTimeout)
	r.opt.RTimeout = types.DecodeInt(r.opt.RTimeout, 0, 3, r.opt.RTimeout)
	r.opt.WTimeout = types.DecodeInt(r.opt.WTimeout, 0, 3, r.opt.WTimeout)
	r.opt.PoolSize = types.DecodeInt(r.opt.PoolSize, 0, 3, r.opt.PoolSize)
	ropts := &redis.UniversalOptions{
		MasterName:   r.opt.MasterName,
		Addrs:        r.opt.Address,
		Password:     r.opt.Password,
		DB:           r.opt.Db,
		DialTimeout:  time.Duration(r.opt.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(r.opt.RTimeout) * time.Second,
		WriteTimeout: time.Duration(r.opt.WTimeout) * time.Second,
		PoolSize:     r.opt.PoolSize,
	}
	r.UniversalClient = redis.NewUniversalClient(ropts)
	_, err = r.UniversalClient.Ping().Result()
	return
}
