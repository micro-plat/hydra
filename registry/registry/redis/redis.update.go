package redis

import (
	"fmt"

	"github.com/micro-plat/lib4go/security/md5"
)

func (r *redisRegistry) Update(path string, data string) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	rpath := joinR(path)
	b, err := r.Exists(path)
	if err != nil {
		return fmt.Errorf("更新节点[%s],检查是否存在异常,err:%+v", path, err)
	}
	if !b {
		return fmt.Errorf("更新节点[%s]不存在", path)
	}

	t, err := r.client.PTTL(rpath).Result()
	if err != nil {
		return err
	}

	res := r.client.Set(rpath, data, t)
	_, err = res.Result()
	if err != nil {
		return fmt.Errorf("更新节点[%s]异常,err:%+v", path, err)
	}

	val := md5.Encrypt(data)
	r.watchMap.Store(path, val)
	return
}
