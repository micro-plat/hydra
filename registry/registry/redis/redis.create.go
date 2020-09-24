package redis

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/components/uuid"
	"github.com/micro-plat/lib4go/security/md5"
)

func (r *redisRegistry) CreatePersistentNode(path string, data string) (err error) {

	if err = r.createNode(path, data, -1); err != nil {
		return
	}
	val := md5.Encrypt(data)
	r.watchMap.Store(path, val)
	return nil
}

func (r *redisRegistry) CreateTempNode(path string, data string) (err error) {

	if err = r.createNode(path, data, LEASE_TTL*time.Second); err != nil {
		return fmt.Errorf("创建临时节点异常,err:%+v", err)
	}
	r.leases.Store(path, data)
	return
}

func (r *redisRegistry) CreateSeqNode(path string, data string) (rpath string, err error) {
	if r.done {
		err = ErrClientConnClosing
		return
	}
	nid := uuid.Get(r.client.ClientID().String())
	rpath = fmt.Sprintf("%s%d", path, nid)
	if err = r.createNode(rpath, data, LEASE_TTL*time.Second); err != nil {
		return "", fmt.Errorf("创建临时seq节点异常,err:%+v", err)
	}

	r.leases.Store(rpath, data)
	return rpath, nil
}

func (r *redisRegistry) createNode(path, data string, timeout time.Duration) (err error) {

	if r.done {
		return ErrClientConnClosing
	}

	b, err := r.Exists(path)
	if err != nil {
		return err
	}
	if b {
		return nil
	}

	ch := make(chan error, 1)
	go func(ch chan error) {
		rpath := joinR(path)
		_, err = r.client.Set(rpath, data, timeout).Result()
		if err != nil {
			ch <- fmt.Errorf("创建[%s]异常,err:%+v", path, err)
			return
		}
		ch <- nil
	}(ch)

	// 使用计时器判断创建节点是否超时
	select {
	case <-time.After(r.options.Timeout):
		err = fmt.Errorf("create node : %s timeout", path)
		return
	case err = <-ch:
		return
	}
}
