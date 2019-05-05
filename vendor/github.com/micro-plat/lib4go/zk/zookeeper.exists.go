package zk

import (
	"fmt"
	"time"
)

//ExistsAny 是否有一个路径已经存在
func (client *ZookeeperClient) ExistsAny(paths ...string) (b bool, path string, err error) {
	for _, path = range paths {
		if b, err = client.Exists(path); err != nil || b {
			return
		}
	}
	return
}

type existsType struct {
	b   bool
	err error
}

//Exists 检查路径是否存在
func (client *ZookeeperClient) Exists(path string) (b bool, err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	// 启动一个协程，判断节点是否存在
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		if client.conn == nil {
			return
		}
		b, _, err = client.conn.Exists(path)
		ch <- existsType{b: b, err: err}
	}(ch)

	select {
	case <-time.After(TIMEOUT):
		if client.done {
			return false, ErrClientConnClosing
		}
		err = fmt.Errorf("judgment node : %s exists timeout", path)
		return
	case data := <-ch:
		if client.done {
			return false, ErrClientConnClosing
		}
		err = data.(existsType).err
		if err != nil {
			return false, err
		}
		b = data.(existsType).b
		return b, nil
	}
}
