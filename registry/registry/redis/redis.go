package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/global"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/redis/internal"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

type Redis struct {
	closeCh       chan struct{}
	seqPath       string
	once          sync.Once
	maxSeq        int64
	maxExpiration time.Duration
	tmpExpiration time.Duration
	checkTicker   time.Duration
	tmpNodes      cmap.ConcurrentMap
	client        *internal.Client
}

//NewRedisBy 构建redis注册中心
func NewRedisBy(master string, pwd string, addrs []string, db int, poolSize int) (*Redis, error) {
	return NewRedis(&internal.ClientConf{
		MasterName: master,
		Password:   pwd,
		Address:    addrs,
		Db:         db,
		PoolSize:   poolSize,
	})
}

//NewRedis 构建redis注册中心
func NewRedis(c *internal.ClientConf) (*Redis, error) {
	rds, err := internal.NewClientByConf(c)
	if err != nil {
		return nil, err
	}
	redis := &Redis{
		client:        rds,
		maxExpiration: time.Hour * 24 * 365 * 10,
		tmpExpiration: time.Second * 5,
		checkTicker:   time.Second * 2,
		maxSeq:        9999999999,
		tmpNodes:      cmap.New(4),
		closeCh:       make(chan struct{}),
		seqPath:       swapKey(fmt.Sprintf("hydra/%s/seq", global.Version)),
	}
	go redis.keepalive()
	return redis, nil
}

//Exists 检查节点是否存在
func (r *Redis) Exists(path string) (bool, error) {
	key := swapKey(path)
	e, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if err == nil && e == 1 {
		return true, nil
	}
	//npaths, err := r.client.Keys(key + ":*").Result()
	exists, err := r.client.ExistsChildren(key + ":*")
	if err != nil {
		return false, err
	}
	return exists, err
}

//Close 关闭当前服务
func (r *Redis) Close() error {
	r.once.Do(func() {
		close(r.closeCh)
		r.client.Close()
		r.tmpNodes.Clear()
	})
	return nil
}

//redisFactory 基于redis的注册中心
type redisFactory struct {
	opts *r.Options
}

//Build 根据配置生成文件系统注册中心
func (z *redisFactory) Create(opts ...r.Option) (r.IRegistry, error) {
	for i := range opts {
		opts[i](z.opts)
	}
	conf := &internal.ClientConf{
		Address:    z.opts.Addrs,
		Password:   z.opts.Auth.Password,
		MasterName: z.opts.Auth.Username,
		Db:         types.GetInt(z.opts.Metadata["db"]),
		PoolSize:   types.GetMax(types.GetInt(z.opts.Metadata["pool_size"]), 10),
	}
	return NewRedis(conf)
}

func init() {
	r.Register(r.Redis, &redisFactory{
		opts: &r.Options{},
	})
}
