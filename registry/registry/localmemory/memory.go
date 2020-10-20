package localmemory

import (
	"sync"

	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
)

//Local 本地内存作为注册中心
var Local r.IRegistry = NewLocalMemory()

type localMemory struct {
	closeCh        chan struct{}
	nodes          map[string]*value
	seqValue       int32
	lock           sync.RWMutex
	vlock          sync.Mutex
	clock          sync.Mutex
	valueWatchs    map[string]chan registry.ValueWatcher
	childrenWatchs map[string]chan registry.ChildrenWatcher
}

func NewLocalMemory() *localMemory {
	return &localMemory{
		seqValue:       10000,
		closeCh:        make(chan struct{}),
		nodes:          make(map[string]*value),
		valueWatchs:    make(map[string]chan registry.ValueWatcher),
		childrenWatchs: make(map[string]chan registry.ChildrenWatcher),
	}
}

func (l *localMemory) Exists(path string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	np := r.Join(path)
	if _, ok := l.nodes[np]; ok {
		return true, nil
	}
	return false, nil
}

func (l *localMemory) Close() error {
	return nil
}

//zkRegistry 基于zookeeper的注册中心
type lmFactory struct {
	opts *r.Options
}

//Build 根据配置生成文件系统注册中心
func (z *lmFactory) Create(opts ...r.Option) (r.IRegistry, error) {
	for i := range opts {
		opts[i](z.opts)
	}
	return Local, nil
}

func init() {
	r.Register(r.LocalMemory, &lmFactory{
		opts: &r.Options{},
	})
}
