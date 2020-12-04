package redis

import (
	"fmt"
	"time"
)

//CreatePersistentNode 创建永久节点
func (r *Redis) CreatePersistentNode(path string, data string) (err error) {
	key := swapKey(path)
	value := newValue(data, false)
	_, err = r.client.Set(key, value.String(), r.maxExpiration).Result()
	if err != nil {
		return err
	}
	r.notifyParentChange(key, value.Version)
	return nil
}

//CreateTempNode 创建临时节点
func (r *Redis) CreateTempNode(path string, data string) (err error) {
	key := swapKey(path)
	value := newValue(data, true)
	_, err = r.client.Set(key, value.String(), r.tmpExpiration).Result()
	if err != nil {
		return err
	}
	r.tmpNodes.Set(key, 0)
	r.notifyParentChange(key, value.Version)
	return nil
}

//CreateSeqNode 创建序列节点
func (r *Redis) CreateSeqNode(path string, data string) (rpath string, err error) {

	nid, err := r.getSeq()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s_%d", swapKey(path), nid)
	value := newValue(data, true)
	_, err = r.client.Set(key, value.String(), r.tmpExpiration).Result()
	if err != nil {
		return swapPath(key), err
	}
	r.tmpNodes.Set(key, 0)
	r.notifyParentChange(key, value.Version)
	return swapPath(key), nil
}

//getSeq 处理seq最大值问题
func (r *Redis) getSeq() (int64, error) {
	//获取seq编号
	nid, err := r.client.Incr(r.seqPath).Result()
	if err != nil {
		return 0, err
	}
	if nid >= r.maxSeq {
		r.client.DecrBy(r.seqPath, r.maxSeq)
		return r.getSeq()
	}
	return nid, nil

}
func (r *Redis) keepalive() {
	tk := time.NewTicker(r.checkTicker)
	for {
		select {
		case <-r.closeCh:
			r.tmpNodes.RemoveIterCb(func(key string, v interface{}) bool {
				r.Delete(key)
				return true
			})
			return
		case <-tk.C:
			items := r.tmpNodes.Items()
			for k := range items {
				if ok, err := r.Exists(k); ok && err == nil {
					r.client.Expire(k, r.tmpExpiration).Result()
				}
			}
		}
	}
}
