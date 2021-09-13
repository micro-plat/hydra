package mysql

import (
	"encoding/json"
	"path/filepath"

	"github.com/micro-plat/lib4go/registry"
)

type valueWatcher struct {
	watcher chan registry.ValueWatcher
}

type childrenWatcher struct {
	watcher chan registry.ChildrenWatcher
}

//WatchValue 监控值变化
func (r *Mysql) WatchValue(path string) (data chan registry.ValueWatcher, err error) {

	if watcher, ok := r.valueWatcherMaps[path]; ok {
		return watcher.watcher, nil
	}
	watcher := make(chan registry.ValueWatcher, 1)

	r.valueWatcherMaps[path] = &valueWatcher{
		watcher: watcher,
	}

	msgChan := r.client.Subscribe(path).ChannelValue()

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
				watcher <- &valueEntity{path: path, version: nv.Version, Value: []byte(nv.Data), Err: err}
			}
		}

	}()
	return watcher, nil
}

//notifyValueChange 通知订阅者值已发生变化
func (r *Mysql) notifyValueChange(path string, value *value) {
	r.client.PublishValue(path, value.String())
}

//WatchChildren 监控子节点变化，保存订阅者信息
func (r *Mysql) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	if watcher, ok := r.childrenWatcherMaps[path]; ok {
		return watcher.watcher, nil
	}
	watcher := make(chan registry.ChildrenWatcher, 1)

	r.childrenWatcherMaps[path] = &childrenWatcher{
		watcher: watcher,
	}

	msgChan := r.client.Subscribe(path).ChannelChild()

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
func (r *Mysql) notifyParentChange(path string, version int32) {
	parent := filepath.Dir(path)
	_, ok := r.childrenWatcherMaps[parent]
	if !ok {
		return
	}
	bytes, _ := json.Marshal(map[string]string{
		"path": parent,
	})
	r.client.PublishChild(parent, string(bytes))
}
