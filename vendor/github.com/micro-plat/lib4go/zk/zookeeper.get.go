package zk

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/encoding"
)

type getValueType struct {
	data    []byte
	version int32
	err     error
}

//GetValue 获取节点的值
func (client *ZookeeperClient) GetValue(path string) (value []byte, version int32, err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	// 起一个协程，获取节点的值
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, stat, err := client.conn.Get(path)
		ch <- getValueType{data: data, err: err, version: stat.Version}
	}(ch)

	select {
	case <-time.After(TIMEOUT):
		err = fmt.Errorf("get node:%s value timeout", path)
		return
	case data := <-ch:
		if client.done {
			err = ErrClientConnClosing
			return
		}
		err = data.(getValueType).err
		if err != nil {
			err = fmt.Errorf("get node:%s error(err:%v)", path, err)
			return
		}
		value, err = encoding.DecodeBytes(data.(getValueType).data, "gbk")
		if err != nil {
			err = fmt.Errorf("get node 编码转换失败:%s error(err:%v)", path, err)
			return
		}
		version = data.(getValueType).version
		return
	}
}

type getChildrenType struct {
	data    []string
	version int32
	err     error
}

//GetChildren 获取节点下的子节点
func (client *ZookeeperClient) GetChildren(path string) (paths []string, version int32, err error) {
	if !client.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if client.done {
		err = ErrClientConnClosing
		return
	}
	if b, err := client.Exists(path); !b || err != nil {
		return nil, 0, fmt.Errorf("node(%s) is not exist", path)
	}

	// 起一个协程，获取子节点
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, stat, err := client.conn.Children(path)
		ch <- getChildrenType{data: data, err: err, version: stat.Version}
	}(ch)

	// 使用定时器判断获取子节点是否超时
	select {
	case <-time.After(TIMEOUT):
		err = fmt.Errorf("get node(%s) children timeout ", path)
		return
	case data := <-ch:
		if client.done {
			err = ErrClientConnClosing
			return
		}
		paths = data.(getChildrenType).data
		version = data.(getChildrenType).version
		err = data.(getChildrenType).err
		if err != nil {
			err = fmt.Errorf("get node(%s) children error(err:%v)", path, err)
		}
		return
	}
}
