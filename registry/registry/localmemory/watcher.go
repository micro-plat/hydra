package localmemory

import "github.com/micro-plat/lib4go/registry"

func (l *localMemory) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	l.vlock.Lock()
	defer l.vlock.Unlock()
	if d, ok := l.valueWatchs[path]; ok {
		return d, nil
	}

	watcher := make(chan registry.ValueWatcher)
	l.valueWatchs[path] = watcher
	return watcher, nil

}

func (l *localMemory) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {
	l.clock.Lock()
	defer l.clock.Unlock()
	if d, ok := l.childrenWatchs[path]; ok {
		return d, nil
	}
	watcher := make(chan registry.ChildrenWatcher)

	l.childrenWatchs[path] = watcher
	return watcher, nil
}

func (l *localMemory) notifyValueChange(path string) {
	l.vlock.Lock()
	defer l.vlock.Unlock()
	for k, v := range l.valueWatchs {
		if k == path {
			v <- &valueEntity{path: path}
		}
		break
	}
	delete(l.valueWatchs, path)
}
func (l *localMemory) notifyParentChange(path string) {
	// path, _, err := l.GetParent(path)
	// if err != nil {
	// 	return
	// }
	// l.vlock.Lock()
	// defer l.vlock.Unlock()
	// for k, v := range l.childrenWatchs {
	// 	if k == path {
	// 		v <- &valuesEntity{path: path}
	// 	}
	// 	break
	// }
	// delete(l.childrenWatchs, path)
}
