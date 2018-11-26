package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/cache"
	"github.com/micro-plat/lib4go/redis"
)

// redisClient redis配置文件
type redisClient struct {
	servers []string
	client  *redis.Client
}

// New 根据配置文件创建一个redis连接
func New(addrs []string, conf string) (m *redisClient, err error) {
	m = &redisClient{servers: addrs}
	m.client, err = redis.NewClientByJSON(conf)
	if err != nil {
		return
	}
	return
}

// Get 根据key获取redis中的数据
func (c *redisClient) Get(key string) (string, error) {
	data, err := c.client.Get(key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return data, nil
		}
		return "", err
	}
	return data, nil
}

//Decrement 增加变量的值
func (c *redisClient) Decrement(key string, delta int64) (n int64, err error) {
	return c.client.DecrBy(key, delta).Result()
}

//Increment 减少变量的值
func (c *redisClient) Increment(key string, delta int64) (n int64, err error) {
	return c.client.IncrBy(key, delta).Result()
}

//Gets 获取多条数据
func (c *redisClient) Gets(key ...string) (r []string, err error) {
	data, err := c.client.MGet(key...).Result()
	if err != nil {
		return nil, err
	}
	r = make([]string, 0, len(data))
	for _, v := range data {
		if v == nil || v.(string) == "" {
			continue
		}
		r = append(r, v.(string))
	}
	return
}

// Add 添加数据到redis中,如果redis存在，则报错
func (c *redisClient) Add(key string, value string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	i, err := c.client.Exists(key).Result()
	if err != nil {
		return err
	}
	if i == 1 {
		err = fmt.Errorf("key:%s已存在", key)
		return err
	}
	_, err = c.client.Set(key, value, expires).Result()
	return err
}

// Set 更新数据到redis中，没有则添加
func (c *redisClient) Set(key string, value string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	_, err := c.client.Set(key, value, expires).Result()
	return err
}

func (c *redisClient) Delete(key string) error {
	if !strings.Contains(key, "*") {
		_, err := c.client.Del(key).Result()
		if err != nil {
			return fmt.Errorf("%v(%s)", err, key)
		}
		return nil
	}
	_, err := c.client.Eval(`
    local keys=redis.call('KEYS',KEYS[1])
    if (#keys==0) then
        return 0
    end
		return redis.call('DEL',unpack(keys))`, []string{key}).Result()
	return err
}

// Delete 删除redis中的数据
func (c *redisClient) Exists(key string) bool {
	r, err := c.client.Exists(key).Result()
	return err == nil && r == 1
}

// Delay 延长数据在redis中的时间
func (c *redisClient) Delay(key string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	_, err := c.client.Expire(key, expires).Result()
	return err
}
func (c *redisClient) Close() error {
	return c.client.Close()
}

type redisResolver struct {
}

func (s *redisResolver) Resolve(address []string, conf string) (cache.ICache, error) {
	return New(address, conf)
}
func init() {
	cache.Register("redis", &redisResolver{})
}
