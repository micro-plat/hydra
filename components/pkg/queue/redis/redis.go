package redis

import (
	rds "github.com/go-redis/redis"
	"github.com/micro-plat/lib4go/queue"
	"github.com/micro-plat/lib4go/redis"
)

// redisClient memcache配置文件
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

// Push 向存于 key 的列表的尾部插入所有指定的值
func (c *redisClient) Push(key string, value string) error {
	_, err := c.client.RPush(key, value).Result()
	return err
}

// Pop 移除并且返回 key 对应的 list 的第一个元素。
func (c *redisClient) Pop(key string) (string, error) {
	r, err := c.client.LPop(key).Result()
	if err != nil && err == rds.Nil {
		return "", queue.Nil
	}
	return r, err
}

// Pop 移除并且返回 key 对应的 list 的第一个元素。
func (c *redisClient) Count(key string) (int64, error) {
	return c.client.LLen(key).Result()
}

// Close 释放资源
func (c *redisClient) Close() error {
	return c.client.Close()
}

type redisResolver struct {
}

func (s *redisResolver) Resolve(address []string, conf string) (queue.IQueue, error) {
	return New(address, conf)
}
func init() {
	queue.Register("redis", &redisResolver{})
}
