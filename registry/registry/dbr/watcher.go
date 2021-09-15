package dbr

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
func (r *DBR) WatchValue(path string) (data chan registry.ValueWatcher, err error) {

	if watcher, ok := r.valueWatcherMaps[path]; ok {
		return watcher.watcher, nil
	}

	watcher := make(chan registry.ValueWatcher, 1)
	r.valueWatcherMaps[path] = &valueWatcher{
		watcher: watcher,
	}

	// watcher <- &valueEntity{path: path, version: nv.Version, Value: []byte(nv.Data), Err: err}

	return watcher, nil
}

//notifyValueChange 通知订阅者值已发生变化
func (r *DBR) notifyValueChange(path string, value string) {
	r.client.Publish(path, value)
}

//WatchChildren 监控子节点变化，保存订阅者信息
func (r *DBR) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	if watcher, ok := r.childrenWatcherMaps[path]; ok {
		return watcher.watcher, nil
	}
	watcher := make(chan registry.ChildrenWatcher, 1)

	r.childrenWatcherMaps[path] = &childrenWatcher{
		watcher: watcher,
	}

	// watcher <- &childrenEntity{path: path, version: ver, children: children, Err: err}

	return watcher, nil
}

//notifyParentChange 通知订阅者值已发生变化
func (r *DBR) notifyParentChange(path string, version int32) {
	parent := filepath.Dir(path)
	_, ok := r.childrenWatcherMaps[parent]
	if !ok {
		return
	}
	bytes, _ := json.Marshal(map[string]string{
		"path": parent,
	})
	r.client.Publish(parent, string(bytes))
}
