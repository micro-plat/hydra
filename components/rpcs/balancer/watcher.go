package balancer

import (
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry"
	r "github.com/micro-plat/lib4go/registry"

	"sort"

	"fmt"

	"sync"

	"google.golang.org/grpc/naming"
)

// Watcher is the implementaion of grpc.naming.Watcher
type Watcher struct {
	client        registry.IRegistry
	isInitialized bool
	caches        map[string]bool
	plat          string
	server        string
	path          string
	service       string
	sortPrefix    string
	closeCh       chan struct{}
	lastErr       error
	once          sync.Once
}

// Close do nothing
func (w *Watcher) Close() {
	w.once.Do(func() {
		close(w.closeCh)
	})
}

// Next 监控服务器地址变化,监控发生异常时移除所有服务,否则等待服务器地址变化
func (w *Watcher) Next() ([]*naming.Update, error) {
start:
	w.lastErr = nil
	if !w.isInitialized {
		path, err := w.initialize()
		if err != nil {
			time.Sleep(time.Second)
			return w.getUpdates([]string{}), nil
		}
		w.path = path
		resp, _, err := w.client.GetChildren(w.path)
		if err == nil {
			w.isInitialized = true
			addrs := w.extractAddrs(resp)
			return w.getUpdates(addrs), nil
		}
	}

	// generate etcd/zk Watcher
	watcherCh, err := w.client.WatchChildren(w.path)
	if err != nil {
		return nil, fmt.Errorf("rpc.client.未找到服务:%s(err:%v)", w.path, err)
	}
	var watcher r.ChildrenWatcher
	select {
	case watcher = <-watcherCh:
	case <-w.closeCh:
		return w.getUpdates([]string{}), w.lastErr
	}
	if err = watcher.GetError(); err != nil {
		return nil, err
	}
	chilren, _ := watcher.GetValue()
	if len(chilren) == 0 {
		w.isInitialized = false
		goto start
	}
	addrs := w.extractAddrs(chilren)
	return w.getUpdates(addrs), nil
}
func (w *Watcher) getUpdates(addrs []string) (updates []*naming.Update) {
	newCache := make(map[string]bool)
	for i := 0; i < len(addrs); i++ {
		newCache[addrs[i]] = true
		if _, ok := w.caches[addrs[i]]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Add, Addr: addrs[i]})
		} else {
			w.caches[addrs[i]] = false
		}
	}
	for i, v := range w.caches {
		if v {
			updates = append(updates, &naming.Update{Op: naming.Delete, Addr: i})
		}
	}
	w.caches = newCache
	return
}
func (w *Watcher) extractAddrs(resp []string) []string {
	addrs := make([]string, 0, len(resp))
	for _, v := range resp {
		item := strings.SplitN(v, "_", 2)
		addrs = append(addrs, item[0])
	}
	if w.sortPrefix != "" {
		sort.Slice(addrs, func(i, j int) bool {
			return strings.HasPrefix(addrs[i], w.sortPrefix)
		})
	}
	return addrs
}

func (w *Watcher) initialize() (string, error) {
	rpath := fmt.Sprintf("/%s/services/rpc/%s%s/providers", w.plat, w.server, w.service)
	b, err := w.client.Exists(rpath)
	if err != nil {
		return "", err
	}
	if b {
		clds, _, err := w.client.GetChildren(rpath)
		if err != nil {
			return "", err
		}
		if len(clds) > 0 {
			return rpath, nil
		}
	}
	root := fmt.Sprintf("/%s/services/rpc/%s", w.plat, w.server)
	items := strings.Split(strings.Trim(w.service, "/"), "/")
	list, err := w.findRealPath([]string{root}, items...)
	if err != nil {
		return "", err
	}
	for _, li := range list {
		path := registry.Join(li, "providers")
		clds, _, err := w.client.GetChildren(path)
		if err != nil {
			return "", err
		}
		if len(clds) > 0 {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到服务提供程序:%s", rpath)

}

func (w *Watcher) findRealPath(roots []string, items ...string) ([]string, error) {
	if len(items) == 0 {
		return roots, nil
	}
	rmatch := make([]string, 0, 1)
	for _, root := range roots {
		paths, _, err := w.client.GetChildren(root)
		if err != nil {
			return nil, err
		}
		for _, p := range paths {
			if p == items[0] || strings.HasPrefix(p, ":") {
				rmatch = append(rmatch, registry.Join(root, p))
			}
		}
	}
	if len(rmatch) == 0 {
		return nil, fmt.Errorf("未找到匹配的路径：%v %s", roots, items[0])
	}
	return w.findRealPath(rmatch, items[1:]...)
}
