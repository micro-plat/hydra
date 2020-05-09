package zk

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/encoding"
)

// Update 更新一个节点的值，如果存在则更新，如果不存在则报错
func (client *ZookeeperClient) Update(path string, data string, version int32) (err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	// 判断节点是否存在
	if b, err := client.Exists(path); !b || err != nil {
		return fmt.Errorf("update node %s fail(node is not exists : %t, err : %v)", path, b, err)
	}

	// 启动一个协程，更新节点
	ch := make(chan error, 1)
	go func(ch chan error) {
		buff, err := encoding.Encode(data, "gbk")
		if err != nil {
			ch <- err
			return
		}
		_, err = client.conn.Set(path, buff, version)
		ch <- err
	}(ch)

	// 启动一个计时器，判断更新节点是否超时
	select {
	case <-time.After(TIMEOUT):
		err = fmt.Errorf("update node %s timeout", path)
		return
	case err = <-ch:
		return err
	}
}
