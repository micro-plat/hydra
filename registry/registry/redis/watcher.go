package redis

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/registry"
	"github.com/micro-plat/lib4go/utility"
)

//WatchValue 监控值变化
func (r *Redis) WatchValue(path string) (data chan registry.ValueWatcher, err error) {

	//为当前监听分配编号
	id := utility.GetGUID()
	key := swapKey(path, "watch") //保存所有watcher编号
	_, err = r.client.HSet(key, id, "watch.value").Result()
	if err != nil {
		return nil, err
	}

	//等待值变化
	nwatch := swapKey(key, id)
	watcher := make(chan registry.ValueWatcher, 1)
	go func() {

		defer r.Delete(nwatch)
		defer r.client.HDel(key, id).Result()

		//拉取数据
		var v []string
		var err error
	LOOP:
		for {
			select {
			case <-r.closeCh: //系统退出
				return
			default:
				v, err = r.client.BLPop(time.Second, nwatch).Result() //获取数据
				if err == nil {                                       //有数据，返回
					break LOOP
				}
				if err != nil && err.Error() != "redis: nil" { //有错误，退出
					watcher <- &valueEntity{Err: err}
					return
				}
			}
		}

		//数据错误
		if len(v) == 0 {
			watcher <- &valueEntity{Err: fmt.Errorf("未收到任何数据")}
			return
		}

		//数据错误
		nv, err := newValueByJSON(v[len(v)-1])
		if err != nil {
			watcher <- &valueEntity{Err: err}
			return
		}
		//通知变更
		watcher <- &valueEntity{path: path, version: nv.Version, Value: nv.Data, Err: err}

	}()
	return watcher, nil

}

//notifyValueChange 通知订阅者值已发生变化
func (r *Redis) notifyValueChange(path string, value *value) {
	key := swapKey(path, "watch") //保存所有watcher编号
	m, err := r.client.HGetAll(key).Result()
	if err != nil {
		return
	}
	for k := range m {
		nkey := swapKey(key, k)
		r.client.RPush(nkey, value.String()).Result()
	}
}

//WatchChildren 监控子节点变化，保存订阅者信息
func (r *Redis) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	//为当前监听分配编号
	id := utility.GetGUID()
	key := swapKey(path, "watch") //保存所有watcher编号
	r.client.HSet(key, id, "watch.children")

	//等待值变化
	nwatch := swapKey(key, id)
	watcher := make(chan registry.ChildrenWatcher, 1)

	go func() {
		defer r.Delete(nwatch)
		defer r.client.HDel(key, id).Result()

		var v []string
		var err error
	LOOP:
		for {
			select {
			case <-r.closeCh: //系统退出
				return
			default:
				v, err = r.client.BLPop(time.Second, nwatch).Result() //获取数据
				if err == nil {                                       //有数据，返回
					break LOOP
				}
				if err != nil && err.Error() != "redis: nil" { //有错误，退出
					watcher <- &childrenEntity{Err: err}
					return
				}
			}
		}

		//数据错误
		if len(v) == 0 {
			watcher <- &childrenEntity{Err: fmt.Errorf("未收到任何数据")}
			return
		}

		children, ver, err := r.GetChildren(path)
		watcher <- &childrenEntity{path: path, version: ver, children: children, Err: err}
	}()
	return watcher, nil
}

//notifyParentChange 通知订阅者值已发生变化
func (r *Redis) notifyParentChange(path string, version int32) {

	//获取父级
	npath := splitKey(path)
	parent := swapKey(npath[:len(npath)-1]...)

	//获取数据
	buff, err := r.client.Get(parent).Result()
	if err != nil {
		return
	}

	//获取所有通知对象
	key := swapKey(parent, "watch") //保存所有watcher编号
	m, err := r.client.HGetAll(key).Result()
	if err != nil {
		return
	}

	//开始通知
	for k := range m {
		nkey := swapKey(key, k)
		r.client.RPush(nkey, buff).Result()
	}
}
