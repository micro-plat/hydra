package context

import (
	"path/filepath"

	"github.com/micro-plat/hydra/dlock"
)

//NewDLock 基于当前注册中心创建分布式锁
func (c *Context) NewDLock(name string) (lk *dlock.DLock) {
	return dlock.NewLockByRegistry(
		filepath.Join("/", c.GetContainer().GetPlatName(), "locks", name), c.GetContainer().GetRegistry())
}

//NewDLockByRegistry 指定注册中心地址创建分布式锁
func (c *Context) NewDLockByRegistry(name string, registry string) (lk *dlock.DLock, err error) {
	return dlock.NewLock(
		filepath.Join("/", c.GetContainer().GetPlatName(), "locks", name), registry, c.Log)
}
