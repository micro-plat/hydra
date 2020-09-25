package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/components/uuid"
	"github.com/micro-plat/lib4go/security/md5"
)

func (r *redisRegistry) CreatePersistentNode(path string, data string) (err error) {
	if !r.isConnect {
		err = ErrColientCouldNotConnect
		return
	}

	if r.done {
		return ErrClientConnClosing
	}

	b, err := r.Exists(path)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	if path == "/" {
		return nil
	}
	//获取每级目录并检查是否存在，不存在则创建
	paths := r.getPaths(path)
	for i := 0; i < len(paths)-1; i++ {
		b, err := r.Exists(paths[i])
		if err != nil {
			return err
		}
		if b {
			continue
		}
		if err = r.createNode(paths[i], "", -1); err != nil {
			return err
		}
	}

	if err = r.createNode(path, data, -1); err != nil {
		return
	}
	val := md5.Encrypt(data)
	r.watchMap.Store(path, val)
	return nil
}

func (r *redisRegistry) CreateTempNode(path string, data string) (err error) {

	if err = r.CreatePersistentNode(r.GetDir(path), ""); err != nil {
		return
	}

	if err = r.createNode(path, data, LEASE_TTL*time.Second); err != nil {
		return fmt.Errorf("创建临时节点异常,err:%+v", err)
	}
	r.leases.Store(path, data)
	return
}

func (r *redisRegistry) CreateSeqNode(path string, data string) (rpath string, err error) {

	if err = r.CreatePersistentNode(r.GetDir(path), ""); err != nil {
		return
	}

	nid := uuid.Get(r.client.ClientID().String())
	rpath = fmt.Sprintf("%s%d", path, nid)
	if err = r.createNode(rpath, data, LEASE_TTL*time.Second); err != nil {
		return "", fmt.Errorf("创建临时seq节点异常,err:%+v", err)
	}

	r.leases.Store(rpath, data)
	return rpath, nil
}

func (r *redisRegistry) createNode(path, data string, timeout time.Duration) (err error) {

	ch := make(chan error, 1)
	go func(ch chan error) {
		rpath := joinR(path)
		_, err = r.client.Set(rpath, data, timeout).Result()
		if err != nil {
			ch <- fmt.Errorf("创建[%s]异常,err:%+v", path, err)
			return
		}
		ch <- nil
	}(ch)

	// 使用计时器判断创建节点是否超时
	select {
	case <-time.After(r.options.Timeout):
		err = fmt.Errorf("create node : %s timeout", path)
		return
	case err = <-ch:
		return
	}
}

//getPaths 获取当前路径的所有子路径
func (r *redisRegistry) getPaths(path string) []string {
	nodes := strings.Split(path, "/")
	len := len(nodes)
	paths := make([]string, 0, len-1)
	for i := 1; i < len; i++ {
		npath := "/" + strings.Join(nodes[1:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}

//GetDir 获取当前路径的目录
func (r *redisRegistry) GetDir(path string) string {
	paths := r.getPaths(path)
	if len(paths) > 2 {
		return paths[len(paths)-2]
	}
	return "/"
}
