package localmemory

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry"
)

func (l *localMemory) GetValue(path string) (data []byte, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	npath := registry.Format(path)
	if v, ok := l.nodes[npath]; ok {
		return []byte(v.data), v.version, nil
	}

	return nil, 0, fmt.Errorf("节点[%s]不存在", npath)

}
func (l *localMemory) GetChildren(path string) (paths []string, version int32, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	paths = make([]string, 0, 1)
	npath := registry.Format(path)
	exists := make(map[string]string)
	for k := range l.nodes {
		if strings.HasPrefix(k, npath+"/") && len(k) > len(npath) {
			list := strings.Split(strings.Trim(k[len(npath):], "/"), "/")
			name := list[0]
			if _, ok := exists[name]; !ok {
				exists[name] = name
				paths = append(paths, name)
			}
		}
	}
	if v, ok := l.nodes[npath]; ok {
		return paths, v.version, nil
	}
	return paths, 0, nil
}
