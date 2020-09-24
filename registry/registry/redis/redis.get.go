// Package redis provides an redis service registry
package redis

import (
	"fmt"
)

func (r *redisRegistry) GetChildren(path string) (paths []string, version int32, err error) {
	paths = make([]string, 0, 1)
	rpath := joinR(path)
	res := r.client.Keys(fmt.Sprint(rpath, ":*"))
	arry, err := res.Result()
	if err != nil {
		return nil, 0, fmt.Errorf("获取节点[%s]的子节点失败,err:%+v", path, err)
	}

	if arry == nil || len(arry) <= 0 {
		return
	}
	resList := []string{}
	for _, str := range arry {
		resList = append(resList, str[len(fmt.Sprint(rpath, ":")):])
	}
	return resList, 0, nil
}

func (r *redisRegistry) GetValue(path string) (data []byte, version int32, err error) {
	pathKey := joinR(path)
	res := r.client.Get(pathKey)
	b, err := res.Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, 0, fmt.Errorf("节点[%s]不存在", path)
		}
		return nil, 0, fmt.Errorf("获取节点[%s]异常,err:%+v", path, err)
	}
	return b, 0, nil
}
