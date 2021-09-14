package mysql

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/components/dbs"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/client"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/river"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/watcher"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/db"
)

type Mysql struct {
	db                  dbs.IDB
	seqValue            int32
	tmpNodes            cmap.ConcurrentMap
	valueWatcherMaps    map[string]*valueWatcher
	childrenWatcherMaps map[string]*childrenWatcher
	closeCh             chan struct{}
	watcher             *watcher.Watcher
	client              *client.Client
	done                bool
}

func NewMysql(c *xdb.DB, o *r.Options) (*Mysql, error) {
	obj, err := db.NewDB(c.Provider, c.ConnString, c.MaxOpen, c.MaxIdle, c.LifeTime)
	if err != nil {
		return nil, err
	}

	cfg := &river.Config{
		DBConf:        c,
		FlushBulkTime: time.Millisecond * o.FlushTime,
	}

	watcher, err := watcher.NewWatcher(cfg)
	if err != nil {
		return nil, err
	}

	return &Mysql{
		db:                  obj,
		seqValue:            10000,
		tmpNodes:            cmap.New(4),
		watcher:             watcher,
		client:              watcher.GetClient(),
		valueWatcherMaps:    make(map[string]*valueWatcher),
		childrenWatcherMaps: make(map[string]*childrenWatcher),
	}, nil
}

//Close 关闭当前服务
func (r *Mysql) Start() error {
	go func() {
		r.watcher.Watch()
		for {
			select {
			case <-r.closeCh:
				r.watcher.Close()
				return
			}
		}
	}()
	return nil
}

//Close 关闭当前服务
func (r *Mysql) Close() error {
	if r.done {
		return nil
	}
	r.done = true
	close(r.closeCh)
	r.tmpNodes.Clear()
	return nil
}

//mysqlFactory 基于mysql的注册中心
type mysqlFactory struct {
	opts *r.Options
}

//Create 根据配置生成mysql注册中心
func (z *mysqlFactory) Create(opts ...r.Option) (r.IRegistry, error) {
	for i := range opts {
		opts[i](z.opts)
	}

	dbConf := &xdb.DB{
		Provider:   "mysql",
		ConnString: fmt.Sprintf("%s:%s@%s?charset=utf8", z.opts.Auth.Username, z.opts.Auth.Password, z.opts.Addrs[0]),
		MaxOpen:    10,
		MaxIdle:    10,
		LifeTime:   600,
	}

	r, err := NewMysql(dbConf, z.opts)
	if err != nil {
		return nil, err
	}
	r.Start()

	return r, err
}

func init() {
	r.Register(r.Mysql, &mysqlFactory{
		opts: &r.Options{},
	})
}
