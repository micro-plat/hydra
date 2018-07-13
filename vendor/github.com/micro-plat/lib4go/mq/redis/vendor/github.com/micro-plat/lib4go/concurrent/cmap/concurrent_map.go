package cmap

import (
	"encoding/json"
	"sync"
)

//ConcurrentMap  线程安全的MAP，创建SHARD_COUNT个共享map(ConcurrentMapShared), 防止单个锁的瓶颈引用的性能问题
type ConcurrentMap []*ConcurrentMapShared

//ConcurrentMapShared  线程安全的共享map
type ConcurrentMapShared struct {
	items map[string]interface{}
	sync.RWMutex
}

// New 创建ConcurrentMap，线程安全map
func New(count int) ConcurrentMap {
	m := make(ConcurrentMap, count)
	for i := 0; i < count; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

//GetShard  根据KEY获取共享map
func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(len(m))]
}

//MSet 根据map设置值
func (m ConcurrentMap) MSet(data map[string]interface{}) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

//Set  设置指定key的值，如果存在则覆盖
func (m *ConcurrentMap) Set(key string, value interface{}) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

//UpsertCb  Callback to return new element to be inserted into the map
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant
type UpsertCb func(exist bool, valueInMap interface{}, newValue interface{}) interface{}

//Upsert  Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *ConcurrentMap) Upsert(key string, value interface{}, cb UpsertCb) (res interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

// SetIfAbsent  Sets the given value under the specified key if no value was associated with it.
func (m *ConcurrentMap) SetIfAbsent(key string, value interface{}) (bool, interface{}) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	if !ok {
		v = value
		shard.items[key] = v
	}
	shard.Unlock()
	return !ok, v
}

//SetCb  设置回调函数
type SetCb func(input ...interface{}) (interface{}, error)

//SetIfAbsentCb 值不存在是调用回调函数生成对象并添加到ConcurrentMap 中
func (m *ConcurrentMap) SetIfAbsentCb(key string, cb SetCb, input ...interface{}) (ok bool, v interface{}, err error) {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok = shard.items[key]
	if !ok {
		v, err = cb(input...)
		if err != nil {
			shard.Unlock()
			return false, nil, err
		}
		shard.items[key] = v
	}
	shard.Unlock()
	return !ok, v, err
}

//Get  Retrieves an element from map under given key.
func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

//Count  Returns the number of elements within the map.
func (m ConcurrentMap) Count() int {
	count := 0
	rcount := len(m)
	for i := 0; i < rcount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

//Has Looks up an item under specified key
func (m *ConcurrentMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

//Remove  Removes an element from the map.
func (m *ConcurrentMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

//Pop Removes an element from the map and returns it
func (m *ConcurrentMap) Pop(key string) (v interface{}, exists bool) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

//PopAll  Removes an element from the map and returns it
func (m ConcurrentMap) PopAll() (v map[string]interface{}) {
	v = make(map[string]interface{})
	count := len(m)
	ch := make(chan Tuple)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(count)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
					delete(shard.items, key)
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
START:
	for {
		select {
		case tup, ok := <-ch:
			if !ok {
				break START
			}
			v[tup.Key] = tup.Val
		}
	}
	return
}

//IsEmpty Checks if map is empty.
func (m *ConcurrentMap) IsEmpty() bool {
	return m.Count() == 0
}

//Tuple Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val interface{}
}

//Iter Returns an iterator which could be used in a for range loop.
//
// Deprecated: using IterBuffered() will get a better performence
func (m ConcurrentMap) Iter() <-chan Tuple {
	ch := make(chan Tuple)
	count := len(m)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(count)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

//IterBuffered Returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMap) IterBuffered() <-chan Tuple {
	ch := make(chan Tuple, m.Count())
	count := len(m)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(count)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

//Items Returns all items as map[string]interface{}
func (m ConcurrentMap) Items() map[string]interface{} {
	tmp := make(map[string]interface{})

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

//Clear 删除所有元素
func (m ConcurrentMap) Clear() {
	wg := sync.WaitGroup{}
	count := len(m)
	wg.Add(count)
	// Foreach shard.
	for _, shard := range m {
		go func(shard *ConcurrentMapShared) {
			// Foreach key, value pair.
			shard.RLock()
			for key := range shard.items {
				delete(shard.items, key)
			}
			shard.RUnlock()
			wg.Done()
		}(shard)
	}
	wg.Wait()
}

//IterCb Iterator callback,called for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb func(key string, v interface{}) bool

//RemoveCb Iterator callback,返回true从列表中移除key
type RemoveCb func(key string, v interface{}) bool

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *ConcurrentMap) IterCb(fn IterCb) {
	b := false
	for idx := range *m {
		shard := (*m)[idx]
		shard.RLock()
		for key, value := range shard.items {
			if !fn(key, value) {
				b = true
				break
			}
		}
		shard.RUnlock()
		if b {
			break
		}
	}
}

//RemoveIterCb 循环移除
func (m *ConcurrentMap) RemoveIterCb(fn RemoveCb) int {
	count := 0
	for idx := range *m {
		shard := (*m)[idx]
		shard.RLock()
		for key, value := range shard.items {
			if fn(key, value) {
				delete(shard.items, key)
				count++
			}

		}
		shard.RUnlock()
	}
	return count
}

//Keys Return all keys as []string
func (m ConcurrentMap) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(count)
		for _, shard := range m {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = <-ch
	}
	return keys
}

//MarshalJSON Reviles ConcurrentMap "private" variables to json marshal.
func (m ConcurrentMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]interface{})

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// Concurrent map uses Interface{} as its value, therefor JSON Unmarshal
// will probably won't know which to type to unmarshal into, in such case
// we'll end up with a value of type map[string]interface{}, In most cases this isn't
// out value type, this is why we've decided to remove this functionality.
