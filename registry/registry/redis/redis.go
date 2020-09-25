// Package redis provides an redis service registry
package redis

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/micro-plat/hydra/registry"

	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

var (
	ErrColientCouldNotConnect = errors.New("redis: could not connect to the server")
	ErrClientConnClosing      = errors.New("redis: the client connection is closing")
)

//LEASE_TTL 临时节点过期时间
const LEASE_TTL = 30

type redisRegistry struct {
	CloseCh   chan struct{}
	client    redis.UniversalClient
	options   *registry.Options
	lock      sync.RWMutex
	leases    sync.Map
	watchMap  sync.Map
	Log       logger.ILogging
	isConnect bool
	// 是否是手动关闭
	done bool
}

// NewRegistry returns an initialized redis registry
func NewRegistry(opts *registry.Options) (*redisRegistry, error) {
	r := &redisRegistry{
		options: opts,
		CloseCh: make(chan struct{}),
	}

	r.Log = r.options.Logger
	r.options.DialTimeout = types.DecodeInt(r.options.DialTimeout, 0, 3, r.options.DialTimeout)
	r.options.PoolSize = types.DecodeInt(r.options.PoolSize, 0, 3, r.options.PoolSize)
	if r.options.Timeout == 0 {
		r.options.Timeout = 5 * time.Second
	}
	ropts := &redis.UniversalOptions{
		MasterName:   r.options.MasterName,
		Addrs:        r.options.Addrs,
		DB:           r.options.Db,
		DialTimeout:  time.Duration(r.options.DialTimeout) * time.Second,
		ReadTimeout:  r.options.Timeout,
		WriteTimeout: r.options.Timeout,
		PoolSize:     r.options.PoolSize,
	}
	if r.options.Auth != nil {
		ropts.Password = r.options.Auth.Password
	}
	r.client = redis.NewUniversalClient(ropts)
	if _, err := r.client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("redis服务器地址链接异常,err:%+v", err)
	}
	r.isConnect = true
	go r.eventWatch()
	go r.leaseRemain()
	return r, nil
}

func (r *redisRegistry) Close() error {
	close(r.CloseCh)
	r.done = true
	//退出的时候  删除临时节点
	return r.client.Close()
}

func (r *redisRegistry) String() string {
	return "redis"
}

func (r *redisRegistry) leaseRemain() {
	for {
		select {
		case <-r.CloseCh:
			return
		case <-time.After(time.Second * (LEASE_TTL - 3)):
			r.leases.Range(func(key, val interface{}) bool {
				path := key.(string)
				go r.leaseKeepAliveOnce(path)
				return true
			})
		}
	}
}

func (r *redisRegistry) leaseKeepAliveOnce(path string) (err error) {
	b, err := r.Exists(path)
	if err != nil || !b {
		return err
	}

	rpath := joinR(path)
	expires := time.Duration(LEASE_TTL) * time.Second
	if _, err = r.client.Expire(rpath, expires).Result(); err != nil {
		return fmt.Errorf("延期临时节点[%s]异常,err:%+v", path, err)
	}
	return
}

//Join redis路径格式解析
func joinR(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString(":")
	}

	str := strings.ReplaceAll(builder.String(), "/", ":")
	return strings.TrimSuffix(str, ":")
}
