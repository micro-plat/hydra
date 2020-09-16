package dlock

import (
	"testing"

	"github.com/micro-plat/lib4go/logger"
)

func TestTryLock(t *testing.T) {
	lockObj, err := NewLock("taosy-test", "192.168.0.101", logger.New("taosy-log"))
	if err != nil {
		t.Errorf("初始化分布式锁对象异常,err:%+v", err)
		return
	}

	if err = lockObj.TryLock(); err != nil {
		t.Errorf("偿试获取分布式锁异常,err:%+v", err)
		return
	}

	return
}

func TestLock(t *testing.T) {

	return
}

func TestUnlock(t *testing.T) {

	return
}
