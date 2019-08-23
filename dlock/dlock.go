package dlock

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//locks 本地lock缓存
// var locks map[string]*DLock = make(map[string]*DLock)
// var currentLock sync.Mutex

//DLock 分布式锁，基于注册中心实现的分布式锁
type DLock struct {
	name      string
	registry  registry.IRegistry
	path      string
	done      bool
	closeChan chan struct{}
	master    bool
}

//NewLock 构建分布式锁
func NewLock(name string, addr string, l logger.ILogging) (lk *DLock, err error) {
	lk = &DLock{name: name, closeChan: make(chan struct{}, 1)}
	lk.registry, err = registry.NewRegistryWithAddress(addr, l)
	if err != nil {
		return nil, err
	}
	return lk, nil
}

//NewLockByRegistry 根据当前注册中心创建分布式锁
func NewLockByRegistry(name string, r registry.IRegistry) (lk *DLock) {
	lk = &DLock{name: name, registry: r, closeChan: make(chan struct{})}
	return lk
}

//TryLock 偿试获取分布式锁
func (d *DLock) TryLock() (err error) {
	defer func() {
		if err != nil {
			d.registry.Delete(d.path)
		}
	}()
	d.path, err = d.registry.CreateSeqNode(filepath.Join(d.name, "dlock_"),
		fmt.Sprintf(`{"time":%d}`, time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("创建锁%s失败:%v", filepath.Join(d.name, "dlock_"), err)
	}

	cldrs, _, err := d.registry.GetChildren(d.name)
	if err != nil {
		return err
	}
	if isMaster(d.path, d.name, cldrs) {
		return nil
	}
	return fmt.Errorf("未获取到分布式锁")
}

//Lock 以独占方式获取分布式锁
func (d *DLock) Lock(timeout ...time.Duration) (err error) {
	defer func() {
		if err != nil {
			d.registry.Delete(d.path)
		}
	}()
	d.path, err = d.registry.CreateSeqNode(filepath.Join(d.name, "dlock_"), fmt.Sprintf(`{"time":%d}`, time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("创建锁%s失败:%v", filepath.Join(d.name, "dlock_"), err)
	}

	cldrs, _, err := d.registry.GetChildren(d.name)
	if err != nil {
		return err
	}
	if isMaster(d.path, d.name, cldrs) {
		return nil
	}

	//监控子节点变化
	ch, err := d.registry.WatchChildren(d.name)
	if err != nil {
		return err
	}
	deadline := time.Minute
	if len(timeout) > 0 {
		deadline = timeout[0]
	}
	for {
		select {
		case <-time.After(deadline):
			return fmt.Errorf("超时未获取到分布式锁")
		case <-d.closeChan:
			return fmt.Errorf("未获取到分布式锁")
		case cldWatcher := <-ch:
			if cldWatcher.GetError() == nil {
				cldrs, _, _ := d.registry.GetChildren(d.name)
				d.master = isMaster(d.path, d.name, cldrs)
				if d.master {
					return nil
				}
			}
		LOOP:
			ch, err = d.registry.WatchChildren(d.name)
			if err != nil {
				if d.done {
					return fmt.Errorf("未获取到分布式锁")
				}
				time.Sleep(time.Second)
				goto LOOP
			}
		}
	}
}

//Unlock 释放分布式锁
func (d *DLock) Unlock() {
	d.done = true
	close(d.closeChan)
	d.registry.Delete(d.path)
}

func isMaster(path string, root string, cldrs []string) bool {
	if len(cldrs) == 0 {
		return false
	}
	ncldrs := make([]string, 0, len(cldrs))
	for _, v := range cldrs {
		name := strings.Replace(v, root, "", -1)
		ncldrs = append(ncldrs, name)
	}
	sort.Strings(ncldrs)
	return strings.HasSuffix(path, ncldrs[0])
}
