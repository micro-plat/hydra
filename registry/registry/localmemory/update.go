package localmemory

import (
	"fmt"
	"strings"
	"time"

	r "github.com/micro-plat/hydra/registry"
)

func (l *localMemory) Update(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.nodes[r.Join(path)]; ok {
		l.nodes[r.Join(path)] = &value{data: data, version: int32(time.Now().Nanosecond())}
		go l.notifyValueChange(path)
		return nil
	}
	return fmt.Errorf("节点[%s]不存在", path)
}

func (l *localMemory) Delete(path string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	np := r.Join(path)
	for k := range l.nodes {
		if strings.HasPrefix(k, np) {
			delete(l.nodes, k)
			go l.notifyValueChange(path)
			go l.notifyParentChange(path)
		}
	}
	return nil
}
