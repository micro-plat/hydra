package dlock

import "time"

//ILock 分布式鍞
type ILock interface {
	TryLock() (err error)
	Lock(timeout ...time.Duration) (err error)
	Unlock()
}
