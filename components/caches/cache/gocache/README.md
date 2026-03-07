# 自定义内存缓存实现

## 概述
本模块实现了一个简单的内存缓存系统，替代了原有的 `github.com/zkfy/go-cache` 依赖。

## 功能特性
- **线程安全**：使用 `sync.RWMutex` 确保并发读写安全
- **过期时间管理**：支持为缓存项设置过期时间
- **自动清理**：后台线程定期清理过期的缓存项
- **基本操作**：
  - Set：设置缓存项
  - Get：获取缓存项
  - Delete：删除缓存项
  - IncrementInt64：增加整数值
  - DecrementInt64：减少整数值
  - Gets：批量获取缓存项
  - Add：添加缓存项（如果已存在则失败）
  - Exists：检查缓存项是否存在
  - Delay：延长缓存项的过期时间
  - Close：关闭缓存服务

## 实现细节

### 核心结构体

#### cacheItem
```go
type cacheItem struct {
    value      interface{}
    expiration time.Time
}
```
表示单个缓存项，包含值和过期时间。

#### memoryCache
```go
type memoryCache struct {
    items           map[string]cacheItem
    mu              sync.RWMutex
    expiration      time.Duration
    cleanupInterval time.Duration
    stopCleanup     chan bool
}
```
缓存管理器，负责管理所有缓存项和自动清理功能。

### 自动清理机制
- 使用 `time.Ticker` 定期触发清理
- 检查每个缓存项的过期时间，删除已过期的项
- 支持通过 `stopCleanup` 通道停止清理线程

## 使用方法

```go
// 创建缓存实例
cache := newMemoryCache(5*time.Minute, 1*time.Minute) // 默认过期时间5分钟，清理间隔1分钟

// 设置缓存项
cache.Set("key", "value", 10*time.Second)

// 获取缓存项
value, ok := cache.Get("key")
if ok {
    fmt.Println(value)
}
```

## 接口兼容性
本实现完全兼容 `cache.ICache` 和 `cache.ICacheExt` 接口，可以无缝替换原有实现。