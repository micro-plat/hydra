package memcache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/micro-plat/hydra/components/pkgs/cache"
)

// Client memcache配置文件
type Client struct {
	servers []string
	client  *memcache.Client
}

// New 根据配置文件创建一个memcache连接
func New(addrs []string) (m *Client, err error) {
	m = &Client{servers: addrs}
	m.client = memcache.New(addrs...)
	m.client.Timeout = time.Second
	return
}

//GetServers 获取服务器列表
func (c *Client) GetServers() []string {
	return c.servers
}

// Get 根据key获取memcache中的数据
func (c *Client) Get(key string) (string, error) {
	data, err := c.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(data.Value), nil
}

//Decrement 增加变量的值
func (c *Client) Decrement(key string, delta int64) (n int64, err error) {
	value, err := strconv.ParseUint(fmt.Sprintf("%d", delta), 10, 64)
	if err != nil {
		return 0, err
	}
	v, err := c.client.Decrement(key, value)
	if err != nil {
		return
	}
	n, err = strconv.ParseInt(fmt.Sprintf("%d", v), 10, 64)
	if err != nil {
		return 0, err
	}
	return
}

//Increment 减少变量的值
func (c *Client) Increment(key string, delta int64) (n int64, err error) {
	value, err := strconv.ParseUint(fmt.Sprintf("%d", delta), 10, 64)
	if err != nil {
		return 0, err
	}
	v, err := c.client.Increment(key, value)
	if err != nil {
		return
	}
	n, err = strconv.ParseInt(fmt.Sprintf("%d", v), 10, 64)
	if err != nil {
		return 0, err
	}
	return
}

//Gets 获取多条数据
func (c *Client) Gets(key ...string) (r []string, err error) {
	data, err := c.client.GetMulti(key)
	if err != nil {
		return nil, err
	}
	r = make([]string, len(data))
	for _, v := range key {
		r = append(r, string(data[v].Value))
	}
	return
}

// Add 添加数据到memcache中,如果memcache存在，则报错
func (c *Client) Add(key string, value string, expiresAt int) error {
	expires := time.Now().Add(time.Duration(expiresAt) * time.Second).Unix()
	if expiresAt == 0 {
		expires = 0
	}
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expires)}
	return c.client.Add(data)
}

// Set 更新数据到memcache中，没有则添加
func (c *Client) Set(key string, value string, expiresAt int) error {
	expires := time.Now().Add(time.Duration(expiresAt) * time.Second).Unix()
	if expiresAt == 0 {
		expires = 0
	}
	data := &memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expires)}
	err := c.client.Set(data)
	return err
}

//Delete 删除memcache中的数据
func (c *Client) Delete(key string) error {

	return c.client.Delete(key)
}

//Exists 判断指定KEY是否存在
func (c *Client) Exists(key string) bool {
	i, err := c.client.Get(key)
	return err == nil && len(i.Value) > 0
}

// Delay 延长数据在memcache中的时间
func (c *Client) Delay(key string, expiresAt int) error {
	return c.client.Touch(key, int32(expiresAt))
}

// DeleteAll 删除所有缓存数据
func (c *Client) DeleteAll() error {
	return c.client.DeleteAll()
}

//Close 关闭服务
func (c *Client) Close() error {
	return nil
}

type mresolver struct {
}

func (s *mresolver) Resolve(address []string, c string) (cache.ICache, error) {
	return New(address)
}
func init() {
	cache.Register("memcached", &mresolver{})
}
