package gocache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/caches/cache"
	vargocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	"github.com/micro-plat/hydra/global"
)

var _ cache.ICache = (*Client)(nil)
var _ cache.ICacheExt = (*Client)(nil)

// Proto Proto
const Proto = "gocache"

// Client redis配置文件
type Client struct {
	lock    sync.Mutex
	servers []string
	client  *memoryCache
}

// NewByOpts 根据配置文件创建一个gocache连接
func NewByOpts(opts ...vargocache.Option) (m *Client, err error) {
	opt := vargocache.New(opts...)
	return NewByConfig(opt)
}

// NewByConfig 根据配置文件创建一个gocache连接
func NewByConfig(opt *vargocache.GoCache) (m *Client, err error) {
	m = &Client{}
	m.client = newMemoryCache(opt.Expiration, opt.CleanupInterval)
	m.servers = []string{
		global.LocalIP(),
	}
	return
}

// GetServers 获取服务器列表
func (c *Client) GetServers() []string {
	return c.servers
}

// GetProto 获取服务类型
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

// Decrement 增加变量的值
func (c *Client) Decrement(key string, delta int64) (n int64, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	err = c.procInt64val(key)
	if err != nil {
		return
	}
	return c.client.DecrementInt64(key, delta)
}

// Increment 减少变量的值
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

// Gets 获取多条数据
func (c *Client) Gets(key ...string) (r []string, err error) {
	r = make([]string, 0, len(key))
	for _, k := range key {
		v, ok := c.client.Get(k)
		if !ok {
			r = append(r, "")
		} else {
			r = append(r, fmt.Sprint(v))
		}
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

// Delete 删除指定key的缓存
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

// Close 关闭缓存服务
func (c *Client) Close() error {
	// 如果需要清理资源，可以在这里添加
	return nil
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// memoryCache 自定义内存缓存
type memoryCache struct {
	items           map[string]cacheItem
	mu              sync.RWMutex
	expiration      time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

// newMemoryCache 创建新的缓存实例
func newMemoryCache(expiration, cleanupInterval time.Duration) *memoryCache {
	c := &memoryCache{
		items:           make(map[string]cacheItem),
		expiration:      expiration,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
	}

	// 启动定期清理
	if cleanupInterval > 0 {
		go c.cleanupLoop()
	}

	return c
}

// cleanupLoop 定期清理过期项
func (c *memoryCache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期项
func (c *memoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if item.expiration.IsZero() {
			continue
		}
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// Set 设置缓存项
func (c *memoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if duration > 0 {
		expiration = time.Now().Add(duration)
	}

	c.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}

// SetDefault 设置默认缓存项（永不过期）
func (c *memoryCache) SetDefault(key string, value interface{}) {
	c.Set(key, value, 0)
}

// Get 获取缓存项
func (c *memoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// GetWithExpiration 获取缓存项及其过期时间
func (c *memoryCache) GetWithExpiration(key string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, time.Time{}, false
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return nil, time.Time{}, false
	}

	return item.value, item.expiration, true
}

// Add 添加缓存项（如果已存在则失败）
func (c *memoryCache) Add(key string, value interface{}, duration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		return fmt.Errorf("item %s already exists", key)
	}

	var expiration time.Time
	if duration > 0 {
		expiration = time.Now().Add(duration)
	}

	c.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Delete 删除缓存项
func (c *memoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// IncrementInt64 增加整数
func (c *memoryCache) IncrementInt64(key string, delta int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return 0, fmt.Errorf("item %s not found", key)
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return 0, fmt.Errorf("item %s expired", key)
	}

	val, ok := item.value.(int64)
	if !ok {
		return 0, fmt.Errorf("item %s not an int64", key)
	}

	val += delta
	item.value = val
	c.items[key] = item

	return val, nil
}

// DecrementInt64 减少整数
func (c *memoryCache) DecrementInt64(key string, delta int64) (int64, error) {
	return c.IncrementInt64(key, -delta)
}

type cacheResolver struct {
}

func (s *cacheResolver) Resolve(conf string) (cache.ICache, error) {
	return NewByOpts(vargocache.WithRaw(conf))
}
func init() {
	cache.Register("gocache", &cacheResolver{})
}
