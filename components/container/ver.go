package container

import (
	"fmt"
	"sync"
)

type ver struct {
	current string
	keys    []string
}

type vers struct {
	keys map[string]*ver
	lock sync.Mutex
}

func newVers() *vers {
	return &vers{
		keys: make(map[string]*ver),
	}
}
func (v *vers) Add(tp string, name string, key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	tps := fmt.Sprintf("%s-%s", tp, name)
	if _, ok := v.keys[tps]; !ok {
		v.keys[tps] = &ver{current: key, keys: []string{key}}
		return
	}
	v.keys[tps].current = key

	v.keys[tps].keys = append(v.keys[tps].keys, key)
}
func (v *vers) Remove(f func(key string) bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	for _, ver := range v.keys {
		for i, k := range ver.keys {
			if k != ver.current {
				if f(k) {
					ver.keys = append(ver.keys[0:i], ver.keys[i+1:]...)
				}
			}
		}
	}
}
