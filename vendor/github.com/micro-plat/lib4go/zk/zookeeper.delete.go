package zk

import (
	"fmt"
	"time"
)

//Delete 修改指定节点的值
func (client *ZookeeperClient) Delete(path string) (err error) {
	if !client.isConnect {
		return ErrColientCouldNotConnect
	}

	// 启动一个协程，删除节点
	ch := make(chan error)
	go func(ch chan error) {
		if client.conn != nil {
			ch <- client.conn.Delete(path, -1)
		}
	}(ch)

	// 启动一个计时器，判断删除节点是否超时
	tk := time.NewTicker(TIMEOUT)
	select {
	case _, ok := <-tk.C:
		if ok {
			tk.Stop()
			err = fmt.Errorf("delete node : %s timeout", path)
			return
		}
	case err = <-ch:
		tk.Stop()
		return
	}

	return
}
