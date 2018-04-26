package balancer

import (
	"strings"

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
	w.lastErr = nil
	if !w.isInitialized {
		resp, _, err := w.client.GetChildren(w.service)
		w.isInitialized = true
		if err == nil {
			addrs := w.extractAddrs(resp)
			return w.getUpdates(addrs), nil
		}
	}

	// generate etcd/zk Watcher
	watcherCh, err := w.client.WatchChildren(w.service)
	if err != nil {
		return nil, fmt.Errorf("rpc.client.未找到服务:%s(err:%v)", w.service, err)
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
