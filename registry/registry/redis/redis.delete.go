package redis

import (
	"fmt"
	"time"
)

func (r *redisRegistry) Delete(path string) error {
	if !r.isConnect {
		return ErrColientCouldNotConnect
	}
	if r.done {
		return ErrClientConnClosing
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	// 启动一个协程，删除节点
	ch := make(chan error)
	go func(ch chan error) {
		rpath := joinR(path)
		_, err := r.client.Del(rpath).Result()
		if err != nil {
			ch <- fmt.Errorf("删除节点[%s]异常,err:%+v", path, err)
			return
		}
		ch <- nil
		if _, ok := r.watchMap.Load(path); ok {
			r.watchMap.Delete(path)
		}
	}(ch)

	// 启动一个计时器，判断删除节点是否超时
	select {
	case <-time.After(r.options.Timeout):
		return fmt.Errorf("delete node : %s timeout", path)
	case err := <-ch:
		return err
	}
}
