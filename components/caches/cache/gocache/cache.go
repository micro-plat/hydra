package gocache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/caches/cache"
	vargocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	"github.com/micro-plat/hydra/global"
	gocache "github.com/zkfy/go-cache"
)

//Proto Proto
const Proto = "gocache"

// Client redis配置文件
type Client struct {
	lock    sync.Mutex
	servers []string
	client  *gocache.Cache
}

// NewByOpts 根据配置文件创建一个gocache连接
func NewByOpts(opts ...vargocache.Option) (m *Client, err error) {
	opt := vargocache.New(opts...)
	return NewByConfig(opt)
}

// NewByConfig 根据配置文件创建一个gocache连接
func NewByConfig(opt *vargocache.GoCache) (m *Client, err error) {
	m = &Client{}
	m.client = gocache.New(opt.Expiration, opt.CleanupInterval)
	m.servers = []string{
		global.LocalIP(),
	}
	return
}

//GetServers 获取服务器列表
func (c *Client) GetServers() []string {
	return c.servers
}

//GetProto 获取服务类型
func (c *Client) GetProto() string {
	return Proto
}

// Get 根据key获取redis中的数据
func (c *Client) Get(key string) (string, error) {
	v, ok := c.client.Get(key)
	if !ok {
		return "", nil
	}
	return fmt.Sprint(v), nil
}

//Decrement 增加变量的值
func (c *Client) Decrement(key string, delta int64) (n int64, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	err = c.procInt64val(key)
	if err != nil {
		return
	}
	return c.client.DecrementInt64(key, delta)
}

//Increment 减少变量的值
func (c *Client) Increment(key string, delta int64) (n int64, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	err = c.procInt64val(key)
	if err != nil {
		return
	}
	return c.client.IncrementInt64(key, delta)
}

func (c *Client) procInt64val(key string) (err error) {
	val, exp, ok := c.client.GetWithExpiration(key)
	if !ok {
		c.client.SetDefault(key, int64(0))
	}
	if _, ok := val.(int64); ok {
		return
	}
	newval, err := strconv.ParseInt(fmt.Sprint(val), 10, 64)
	if err != nil {
		err = fmt.Errorf("%s的值：%v,不是有效的数字", key, val)
		return
	}
	c.client.Set(key, newval, exp.Sub(time.Now()))
	return

}

//Gets 获取多条数据
func (c *Client) Gets(key ...string) (r []string, err error) {
	r = make([]string, 0, len(key))
	for _, k := range key {
		v, ok := c.client.Get(k)
		if !ok {
			v = ""
		}
		r = append(r, v.(string))
	}
	return r, nil

}

// Add 添加数据到redis中,如果redis存在，则报错
func (c *Client) Add(key string, value string, expiresAt int) error {
	return c.client.Add(key, value, time.Second*time.Duration(expiresAt))
}

// Set 更新数据到redis中，没有则添加
func (c *Client) Set(key string, value string, expiresAt int) error {
	c.client.Set(key, value, time.Second*time.Duration(expiresAt))
	return nil
}

//Delete 删除指定key的缓存
func (c *Client) Delete(key string) error {
	c.client.Delete(key)
	return nil
}

// Exists 查询key是否存在
func (c *Client) Exists(key string) bool {
	_, ok := c.client.Get(key)
	return ok
}

// Delay 延长数据在redis中的时间
func (c *Client) Delay(key string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	v, ok := c.client.Get(key)
	if !ok {
		return fmt.Errorf("%s值不存在", key)
	}
	c.client.Set(key, v, expires)
	return nil
}

//Close 关闭缓存服务
func (c *Client) Close() error {
	return nil
}

type cacheResolver struct {
}

func (s *cacheResolver) Resolve(conf string) (cache.ICache, error) {
	return NewByOpts(vargocache.WithRaw(conf))
}
func init() {
	cache.Register("gocache", &cacheResolver{})
}
