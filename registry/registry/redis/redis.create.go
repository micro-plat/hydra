package redis

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/components/uuid"
	"github.com/micro-plat/lib4go/security/md5"
)

func (r *redisRegistry) CreatePersistentNode(path string, data string) (err error) {
	rpath := joinR(path)
	res := r.client.Set(rpath, data, -1)
	_, err = res.Result()
	if err != nil {
		return fmt.Errorf("创建持久化节点[%s]异常,err:%+v", path, err)
	}
	val := md5.Encrypt(data)
	r.watchMap.Store(path, val)
	return nil
}

func (r *redisRegistry) CreateTempNode(path string, data string) (err error) {
	rpath := joinR(path)
	res := r.client.Set(rpath, data, LEASE_TTL*time.Second)
	_, err = res.Result()
	if err != nil {
		return fmt.Errorf("创建临时节点[%s]异常,err:%+v", path, err)
	}
	r.leases.Store(path, data)
	return
}

func (r *redisRegistry) CreateSeqNode(path string, data string) (rpath string, err error) {

	nid := uuid.Get(r.client.ClientID().String())
	rpath = fmt.Sprintf("%s%d", path, nid)
	xpath := joinR(rpath)
	res := r.client.Set(xpath, data, LEASE_TTL*time.Second)
	_, err = res.Result()
	if err != nil {
		return "", fmt.Errorf("创建seq节点[%s]异常,err:%+v", path, err)
	}
	r.leases.Store(rpath, data)
	return rpath, nil
}
