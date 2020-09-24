package redis

import (
	"fmt"
)

func (r *redisRegistry) Delete(path string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	rpath := joinR(path)
	_, err := r.client.Del(rpath).Result()
	if err != nil {
		return fmt.Errorf("删除节点[%s]异常,err:%+v", path, err)
	}
	if _, ok := r.watchMap.Load(path); ok {
		r.watchMap.Delete(path)
	}
	return nil
}
