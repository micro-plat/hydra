package uuid

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

// 因为UUID目的是解决分布式下生成唯一id 所以ID中是包含集群和节点编号在内的
// var currentNum int64

const (
	workerBits uint8 = 7  // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点
	numberBits uint8 = 10 // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位数 即每毫秒可生成 2^12-1=4096个唯一ID
	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码，感兴趣的同学可以自己算算试试 -1 ^ (-1 << nodeBits) 这里是不是等于 1023
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift   uint8 = workerBits + numberBits // 时间戳向左的偏移量
	workerShift uint8 = numberBits              // 节点ID向左的偏移量
	// 41位字节作为时间戳数值的话 大约68年就会用完
	epoch   int64 = 1577808000 // 1577808000000 //2020-01-01 0:0:0
	timeout int64 = 3600       // map中的uuid对象失效时间 秒(每一个小时清除一次)
)

type suuid struct {
	tag        string
	currentNum int64
	timestamp  int64
}

var suuidMap = cmap.New(8)

var once sync.Once

func GetSUUID(tag string) *suuid {
	//开启监控清空
	once.Do(func() {
		go clear()
	})
LOOP:
	if !suuidMap.Has(tag) {
		obj := &suuid{tag: tag, currentNum: 0, timestamp: time.Now().Unix()}
		suuidMap.Set(tag, obj)
		return obj
	}
	v, b := suuidMap.Get(tag)
	if !b {
		time.Sleep(50 * time.Millisecond)
		goto LOOP
	}
	return v.(*suuid)
}

func clear() {
	tm := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-global.Def.ClosingNotify():
			suuidMap.Clear()
			break
		case <-tm.C:
			suuidMap.RemoveIterCb(func(key string, v interface{}) bool {
				obj := v.(*suuid)
				now := time.Now().Unix()
				if (now - obj.timestamp) > timeout {
					return true
				}
				return false
			})
		}
	}
}

//Get 获取全局唯一编号每个节点每秒1000个不重复
func (n *suuid) Get() UUID {
	wid := fnv32(n.tag)
	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e9 // 纳秒转秒
	id := atomic.AddInt64(&n.currentNum, 1)
	if atomic.CompareAndSwapInt64(&n.currentNum, numberMax-1, 0) {
		id = atomic.AddInt64(&n.currentNum, 1)
	}
	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	nid := int64((now-epoch)<<timeShift | (int64(wid) % workerMax << workerShift) | id)
	return UUID(nid)
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
