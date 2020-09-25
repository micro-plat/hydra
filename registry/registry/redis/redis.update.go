package redis

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/security/md5"
)

func (r *redisRegistry) Update(path string, data string) (err error) {
	if !r.isConnect {
		return ErrColientCouldNotConnect
	}
	if r.done {
		return ErrClientConnClosing
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	ch := make(chan error, 1)
	go func(ch chan error) {
		rpath := joinR(path)
		b, err := r.Exists(path)
		if err != nil {
			ch <- fmt.Errorf("更新节点[%s],检查是否存在异常,err:%+v", path, err)
			return
		}
		if !b {
			ch <- fmt.Errorf("更新节点[%s]不存在", path)
			return
		}

		t, err := r.client.PTTL(rpath).Result()
		if err != nil {
			ch <- err
			return
		}

		_, err = r.client.Set(rpath, data, t).Result()
		if err != nil {
			ch <- fmt.Errorf("更新节点[%s]异常,err:%+v", path, err)
			return
		}
		ch <- nil
		val := md5.Encrypt(data)
		r.watchMap.Store(path, val)
	}(ch)

	select {
	case <-time.After(r.options.Timeout):
		return fmt.Errorf("update node:%s value timeout", path)
	case err = <-ch:
		return
	}
}
