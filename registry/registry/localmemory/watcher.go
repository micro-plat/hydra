package localmemory

import (
	"fmt"
	"strings"

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
	for k, v := range l.valueWatchs {
		if k == path {
			select {
			case v <- &valueEntity{path: path, version: value.version, Value: []byte(value.data)}:
			default:
			}
			break
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
	path, _, err := l.getParentForNotify(cPath)
	if err != nil { //未找到父节点，无需通知
		return
	}
	//@fix 修复节点变动未进行通知的bug
	for k, v := range l.childrenWatchs {
		if k == path {
			select {
			case v <- &childrenEntity{path: path, children: []string{cPath}, version: version}:
			default:
			}
			break
		}
	}
	delete(l.childrenWatchs, path)
}
func (l *localMemory) getParentForNotify(path string) (string, int32, error) {
	npath := r.Format(path)
	for k := range l.childrenWatchs {
		if strings.HasPrefix(npath, k) && len(npath) > len(k) {
			list := r.Split(npath)
			if r.Join(k, list[len(list)-1]) == npath {
				if v, ok := l.nodes[k]; ok {
					return k, v.version, nil
				}
				return k, 0, nil
			}
		}
	}
	return "", 0, fmt.Errorf("节点[%s]的父节点%w", path, errs.ErrNotExist)

}
