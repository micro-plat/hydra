package localmemory

import (
	"fmt"

	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/registry"
)

func (l *localMemory) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	l.vlock.Lock()
	defer l.vlock.Unlock()
	if d, ok := l.valueWatchs[path]; ok {
		return d, nil
	}

	watcher := make(chan registry.ValueWatcher, 1)
	l.valueWatchs[path] = watcher
	return watcher, nil

}
func (l *localMemory) notifyValueChange(path string, value *value) {
	l.vlock.Lock()
	defer l.vlock.Unlock()

	if v, ok := l.valueWatchs[path]; ok {
		select {
		case v <- &valueEntity{path: path, version: value.version, Value: []byte(value.data)}:
		default:
		}
	}
	delete(l.valueWatchs, path)
}
func (l *localMemory) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	l.clock.Lock()
	defer l.clock.Unlock()
	if d, ok := l.childrenWatchs[path]; ok {
		return d, nil
	}

	watcher := make(chan registry.ChildrenWatcher, 1)

	l.childrenWatchs[path] = watcher
	return watcher, nil
}

func (l *localMemory) notifyParentChange(cPath string, version int32) {
	l.clock.Lock()
	defer l.clock.Unlock()

	paths, err := l.getParentForNotify(cPath)
	if err != nil { //未找到父节点，无需通知
		return
	}

	for _, p := range paths {
		if v, ok := l.childrenWatchs[p]; ok {
			select {
			case v <- &childrenEntity{path: p, children: []string{cPath}, version: 0}:
			default:
			}
			delete(l.childrenWatchs, p)
		}
	}
}

func (l *localMemory) getParentForNotify(path string) ([]string, error) {
	npath := r.Format(path)
	res := []string{}
	list := r.Split(npath)
	if len(list) == 1 {
		return nil, fmt.Errorf("节点[%s]的父节点%w", path, errs.ErrNotExist)
	}
	for i := 0; i < len(list)-1; i++ {
		res = append(res, r.Join(list[:i+1]...))
	}
	return res, nil
}
