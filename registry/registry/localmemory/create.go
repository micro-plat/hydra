package localmemory

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/micro-plat/hydra/registry"
)

func (l *localMemory) CreatePersistentNode(path string, data string) (err error) {
	l.lock.Lock()
	list := strings.Split(registry.Format(path), "/")
	for i := range list {
		path := registry.Join(list[:i]...)
		if _, ok := l.nodes[path]; !ok {
			l.nodes[path] = newValue("{}")
		}
	}
	nvalue := newValue(data)
	l.nodes[registry.Format(path)] = nvalue
	l.lock.Unlock()
	l.notifyParentChange(path, nvalue.version)
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
