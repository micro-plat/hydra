package localmemory

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	r "github.com/micro-plat/hydra/registry"
)

func (l *localMemory) CreatePersistentNode(path string, data string) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	list := strings.Split(r.Join(path), "/")
	for i := range list {
		path := r.Join(list[:i]...)
		if _, ok := l.nodes[path]; !ok {
			l.nodes[path] = &value{data: "{}", version: int32(time.Now().Nanosecond())}
		}
	}
	l.nodes[r.Join(path)] = &value{data: data, version: int32(time.Now().Nanosecond())}
	go l.notifyParentChange(path)
	return nil
}
func (l *localMemory) CreateTempNode(path string, data string) (err error) {
	return l.CreatePersistentNode(path, data)
}
func (l *localMemory) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&l.seqValue, 1)
	rpath = fmt.Sprintf("%s_%d", path, nid)
	err = l.CreatePersistentNode(rpath, data)
	return
}
