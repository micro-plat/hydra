package localmemory

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/micro-plat/hydra/registry"
)

func (r *localMemory) getPaths(path string) []string {
	nodes := strings.Split(strings.Trim(path, "/"), "/")
	len := len(nodes)
	paths := make([]string, 0, len)
	for i := 0; i < len; i++ {
		npath := "/" + strings.Join(nodes[:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}

func (l *localMemory) CreatePersistentNode(path string, data string) (err error) {
	l.lock.Lock()
	path = registry.Format(path)
	paths := l.getPaths(path)
	for _, xpath := range paths {
		if _, ok := l.nodes[xpath]; !ok && xpath != path {
			l.nodes[xpath] = newValue("{}")
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
