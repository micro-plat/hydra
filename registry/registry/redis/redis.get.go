// Package redis provides an redis service registry
package redis

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/security/md5"
)

type getChildrenType struct {
	data    []string
	version int32
	err     error
}

func (r *redisRegistry) GetChildren(path string) (paths []string, version int32, err error) {
	if !r.isConnect {
		return nil, 0, ErrColientCouldNotConnect
	}
	if r.done {
		return nil, 0, ErrClientConnClosing
	}

	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		rpath := joinR(path)
		list := []string{}
		c := uint64(0)
		for {
			res := r.client.Scan(c, fmt.Sprint(rpath, ":*"), 100)
			arry, c1, err := res.Result()
			if err != nil {
				ch <- getChildrenType{data: nil, version: 0, err: fmt.Errorf("获取节点[%s]的子节点失败,err:%+v", path, err)}
				return
			}
			if arry != nil && len(arry) > 0 {
				allArry := make([]string, len(arry)+len(list))
				copy(allArry, list)
				copy(allArry[len(list):], arry)
				list = allArry
			}
			if c1 == 0 {
				break
			}
			c = c1
		}

		if list == nil || len(list) <= 0 {
			ch <- getChildrenType{data: list, version: 0, err: nil}
			return
		}
		resList := []string{}
		for _, str := range list {
			resList = append(resList, str[len(fmt.Sprint(rpath, ":")):])
		}
		ch <- getChildrenType{data: resList, version: 0, err: nil}
	}(ch)

	select {
	case <-time.After(r.options.Timeout):
		return nil, 0, fmt.Errorf("get node:%s value timeout", path)
	case data := <-ch:
		err := data.(getChildrenType).err
		if err != nil {
			return nil, 0, err
		}
		et := data.(getChildrenType)
		return et.data, et.version, et.err
	}
}

type getValueType struct {
	data    []byte
	version int32
	err     error
}

func (r *redisRegistry) GetValue(path string) (data []byte, version int32, err error) {
	if !r.isConnect {
		return nil, 0, ErrColientCouldNotConnect
	}

	if r.done {
		return nil, 0, ErrClientConnClosing
	}
	// fmt.Println("pathpath:", path)
	rpath := joinR(path)
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		val, err := r.client.Get(rpath).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				ch <- getValueType{data: nil, version: 0, err: fmt.Errorf("节点[%s]不存在", path)}
				return
			}
			ch <- getValueType{data: nil, version: 0, err: fmt.Errorf("获取节点[%s]异常,err:%+v", path, err)}
			return
		}

		go func() {
			t, err := r.client.PTTL(rpath).Result()
			if err != nil {
				return
			}
			if -1*time.Millisecond == t {
				r.watchMap.Store(path, md5.Encrypt(val))
			}
		}()

		ch <- getValueType{data: []byte(val), version: 0, err: nil}
	}(ch)

	select {
	case <-time.After(r.options.Timeout):
		return nil, 0, fmt.Errorf("get node:%s value timeout", path)
	case data := <-ch:
		err := data.(getValueType).err
		if err != nil {
			return nil, 0, err
		}
		et := data.(getValueType)
		return et.data, et.version, et.err
	}
}
