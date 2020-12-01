package container

import (
	"sync"
)

//history 历史key记录
type history struct {
	current string
	keys    []string
}

//histories 所有配置的历史key记录
type histories struct {
	records map[string]*history
	lock    sync.Mutex
}

func newHistories() *histories {
	return &histories{
		records: make(map[string]*history),
	}
}

//Add 添加key信息
func (v *histories) Add(group string, key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	his, ok := v.records[group]
	if !ok {
		v.records[group] = &history{current: key, keys: []string{key}}
		return
	}
	his.current = key
	his.keys = append(his.keys, key)
}

//Remove 移除key信息
func (v *histories) Remove(f func(key string) bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	for _, history := range v.records {
		for i, k := range history.keys {
			if k != history.current {
				if f(k) {
					history.keys = append(history.keys[0:i], history.keys[i+1:]...)
				}
			}
		}
	}
}
