package dlock

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//locks 本地lock缓存
var locks map[string]*DLock = make(map[string]*DLock)
var currentLock sync.Mutex

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
	currentLock.Lock()
	defer currentLock.Unlock()
	if lk, ok := locks[name]; ok {
		return lk, nil
	}
	lk = &DLock{name: name, closeChan: make(chan struct{})}
	lk.registry, err = registry.NewRegistryWithAddress(addr, l)
	if err != nil {
		return nil, err
	}
	locks[name] = lk
	return lk, nil
}

//NewLockByRegistry 根据当前注册中心创建分布式锁
func NewLockByRegistry(name string, r registry.IRegistry) (lk *DLock) {
	currentLock.Lock()
	defer currentLock.Unlock()
	if lk, ok := locks[name]; ok {
		return lk
	}
	lk = &DLock{name: name, registry: r, closeChan: make(chan struct{})}
	locks[name] = lk
	return lk
}

//TryLock 偿试获取分布式锁
func (d *DLock) TryLock() (err error) {
	d.done = false
	d.master = false
	var path = d.path
	if path == "" {
		path, err = d.registry.CreateSeqNode(filepath.Join(d.name, "dlock_"),
			fmt.Sprintf(`{"time":%d}`, time.Now().Unix()))
		if err != nil {
			return err
		}
	}
	cldrs, _, err := d.registry.GetChildren(d.name)
	if err != nil {
		return err
	}
	if isMaster(path, d.name, cldrs) {
		d.path = path
		return nil
	}
	d.registry.Delete(path)
	return fmt.Errorf("未获取到分布式锁")
}

//Lock 以独占方式获取分布式锁
func (d *DLock) Lock() (err error) {
	d.done = false
	d.master = false
	if d.path == "" {
		d.path, err = d.registry.CreateSeqNode(filepath.Join(d.name, "dlock_"), fmt.Sprintf(`{"time":%d}`, time.Now().Unix()))
		if err != nil {
			return err
		}
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
	for {
		select {
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
	d.closeChan <- struct{}{}
	d.done = true
	d.path = ""
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
