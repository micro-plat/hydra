package localmemory

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry"
)

func (l *localMemory) Update(path string, data string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.nodes[registry.Format(path)]; ok {
		nvalue := newValue(data)
		l.nodes[registry.Format(path)] = nvalue
		l.notifyValueChange(path, nvalue)
		return nil
	}
	return fmt.Errorf("节点[%s]不存在", path)
}

func (l *localMemory) Delete(path string) error {
	rpath := registry.Format(path)
	b, err := l.Exists(rpath)
	if err != nil {
		return err
	}
	if !b {
		return nil
	}

	_, version, err := l.GetValue(rpath)
	if err != nil {
		return err
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	np := rpath
	for k, nv := range l.nodes {
		if strings.HasPrefix(k, np) {
			delete(l.nodes, k)
			l.notifyValueChange(path, nv)
			l.notifyParentChange(path, version)
		}
	}
	return nil
}
