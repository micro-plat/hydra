package redis

import (
	"fmt"
	"time"
)

type existsType struct {
	b   bool
	err error
}

func (r *redisRegistry) Exists(path string) (bool, error) {
	if !r.isConnect {
		return false, ErrColientCouldNotConnect
	}

	if r.done {
		return false, ErrClientConnClosing
	}

	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		rpath := joinR(path)
		rs, err := r.client.Exists(rpath).Result()
		if err == nil && rs == 1 {
			ch <- existsType{b: true, err: err}
			return
		}
		ch <- existsType{b: false, err: err}
	}(ch)

	select {
	case <-time.After(r.options.Timeout):
		return false, fmt.Errorf("judgment node : %s exists timeout", path)
	case data := <-ch:
		err := data.(existsType).err
		if err != nil {
			return false, err
		}
		et := data.(existsType)
		return et.b, et.err
	}
}
