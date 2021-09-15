package dbr

import (
	"path/filepath"

	"github.com/micro-plat/lib4go/registry"
)

//WatchValue 监控值变化
func (r *DBR) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	data = r.valueWatchers.Watch(path)
	return data, nil
}

//WatchChildren 监控子节点变化，保存订阅者信息
func (r *DBR) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	data = r.childrenWatchers.Watch(path)
	return data, nil
}

//notifyValueChange 通知订阅者值已发生变化
func (r *DBR) notifyValueChange(path string, value string, version int32) {
	r.valueWatchers.Notify(path, version, value)
}

//notifyParentChange 通知订阅者值已发生变化
func (r *DBR) notifyParentChange(path string, version int32) {
	r.childrenWatchers.Notify(filepath.Dir(path), version, path)
}
