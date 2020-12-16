package redis

import (
	"encoding/json"

	"github.com/micro-plat/hydra/registry/registry/redis/internal"
	"github.com/micro-plat/lib4go/registry"
)

type valueWatcher struct {
	watcher chan registry.ValueWatcher
}

type childrenWatcher struct {
	watcher chan registry.ChildrenWatcher
}

//WatchValue 监控值变化
func (r *Redis) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	watchkey := internal.SwapKey(path, "watch")

	if watcher, ok := r.valueWatcherMaps[watchkey]; ok {
		return watcher.watcher, nil
	}
	watcher := make(chan registry.ValueWatcher, 1)

	r.valueWatcherMaps[watchkey] = &valueWatcher{
		watcher: watcher,
	}

	pubsub := r.client.Subscribe(watchkey)
	msgChan := pubsub.Channel()

	//等待值变化

	go func() {
		for {
			select {
			case <-r.closeCh: //系统退出
				return
			case msg := <-msgChan:
				nv, err := newValueByJSON(msg.Payload)
				if err != nil { //有错误，退出
					watcher <- &valueEntity{Err: err}
					continue
				}
				watcher <- &valueEntity{path: path, version: nv.Version, Value: nv.Data, Err: err}
			}
		}

	}()
	return watcher, nil
}

//notifyValueChange 通知订阅者值已发生变化
func (r *Redis) notifyValueChange(path string, value *value) {
	key := internal.SwapKey(path, "watch") //保存所有watcher编号
	r.client.Publish(key, value.String())
}

//WatchChildren 监控子节点变化，保存订阅者信息
func (r *Redis) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {

	watchkey := internal.SwapKey(path, "watch") //保存所有watcher编号

	if watcher, ok := r.childrenWatcherMaps[watchkey]; ok {
		return watcher.watcher, nil
	}
	//等待值变化
	watcher := make(chan registry.ChildrenWatcher, 1)

	r.childrenWatcherMaps[watchkey] = &childrenWatcher{
		watcher: watcher,
	}

	pubsub := r.client.Subscribe(watchkey)
	msgChan := pubsub.Channel()

	go func() {
		for {
			select {
			case <-r.closeCh: //系统退出
				return
			case <-msgChan:
				children, ver, err := r.GetChildren(path)
				if err != nil { //有错误，退出
					watcher <- &childrenEntity{Err: err}
					continue
				}
				watcher <- &childrenEntity{path: path, version: ver, children: children, Err: err}
			}
		}

	}()
	return watcher, nil
}

//notifyParentChange 通知订阅者值已发生变化
func (r *Redis) notifyParentChange(path string, version int32) {
	npath := internal.SplitKey(path)
	parent := internal.SwapKey(npath[:len(npath)-1]...)
	_, ok := r.childrenWatcherMaps[parent]
	if !ok {
		return
	}

	bytes, _ := json.Marshal(map[string]string{
		"path": parent,
	})

	r.client.Publish(parent, string(bytes))
	return
}
