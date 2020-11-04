package dlock

import (
	"testing"
	"time"

	_ "github.com/micro-plat/hydra/registry/registry/zookeeper"
	"github.com/micro-plat/lib4go/logger"
)

//修改点:1.锁的name应该采用 platname/dlock/自定义名称
// 		3.isMaster判断时,由于replace出现空字串,所以一直会返回true

func TestTryLock(t *testing.T) {

	lockObj, err := NewLock("mgrweb", "zk://192.168.0.101:2181", logger.New("taosy-log"))
	if err != nil {
		t.Errorf("初始化分布式锁对象异常,err:%+v", err)
		return
	}

	if err := lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常1,err:%+v", err)
		return
	}
	// ncldrs := make([]string, 0, len(lockObj.Chads))
	// for _, v := range lockObj.Chads {
	// 	name := strings.Replace(v, "mgrweb", "", -1)
	// 	ncldrs = append(ncldrs, name)
	// }

	// sort.Strings(ncldrs)
	// t.Errorf("11111111:Chads:%v,Pathr:%s,bool:%v", ncldrs, lockObj.Pathr, strings.HasSuffix(lockObj.Pathr, ncldrs[0]))
	if err := lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常2,err:%+v", err)
		return
	}

	// ncldrs = make([]string, 0, len(lockObj.Chads))
	// for _, v := range lockObj.Chads {
	// 	name := strings.Replace(v, "mgrweb", "", -1)
	// 	ncldrs = append(ncldrs, name)
	// }

	// sort.Strings(ncldrs)
	// t.Logf("2222222222:Chads:%v,Pathr:%s,bool:%v", ncldrs, lockObj.Pathr, strings.HasSuffix(lockObj.Pathr, ncldrs[0]))
	if err := lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常3,err:%+v", err)
		return
	}
	// ncldrs = make([]string, 0, len(lockObj.Chads))
	// for _, v := range lockObj.Chads {
	// 	name := strings.Replace(v, "mgrweb", "", -1)
	// 	ncldrs = append(ncldrs, name)
	// }

	// sort.Strings(ncldrs)
	// t.Logf("33333333333:Chads:%v,Pathr:%s,bool:%v", ncldrs, lockObj.Pathr, strings.HasSuffix(lockObj.Pathr, ncldrs[0]))
	return
}

//修改点:1.锁的name应该采用 platname/dlock/自定义名称
// 		3.isMaster判断时,由于replace出现空字串,所以一直会返回true
//		4.锁获取等待默认时间过长;
//		5.LOOP 中WatchChildren出错,可能会陷入死循环
//		6.WatchChildren如参路径有错误
//		7.CreateSeqNode如参路径有错误

func TestLock(t *testing.T) {
	lockObj, err := NewLock("mgrweb", "zk://192.168.0.101:2181", logger.New("taosy-log"))
	if err != nil {
		t.Errorf("初始化分布式锁对象异常,err:%+v", err)
		return
	}

	lockObj1, err := NewLock("mgrweb", "zk://192.168.0.101:2181", logger.New("taosy-log"))
	if err != nil {
		t.Errorf("初始化分布式锁对象异常,err:%+v", err)
		return
	}

	if err := lockObj.Lock(); err != nil {
		t.Errorf("偿试获取分布式锁异常1,err:%+v", err)
		return
	}

	go func() {
		time.Sleep(5 * time.Second)
		lockObj.Unlock()
	}()

	if err := lockObj1.Lock(6 * time.Second); err != nil {
		t.Errorf("偿试获取分布式锁异常1,err:%+v", err)
		return
	}

	go func() {
		time.Sleep(10 * time.Second)
		lockObj1.Unlock()
	}()

	if err := lockObj.Lock(6 * time.Second); err == nil {
		t.Error("不应该拿到锁")
		return
	}

	return
}

//不太明白closeChan的作用
func TestUnlock(t *testing.T) {
	lockObj, err := NewLock("mgrweb", "zk://192.168.0.101:2181", logger.New("taosy-log"))
	if err != nil {
		t.Errorf("初始化分布式锁对象异常,err:%+v", err)
		return
	}

	if err := lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常3,err:%+v", err)
		return
	}

	if err := lockObj.TryLock(); err == nil {
		t.Error("不应该取到锁", err)
		return
	}

	lockObj.Unlock()
	if err := lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常3,err:%+v", err)
		return
	}

	TestLock(t)
	return
}
