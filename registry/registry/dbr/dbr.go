package dbr

import (
	"fmt"

	"github.com/micro-plat/hydra/components/dbs"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/db"
)

type DBR struct {
	db               dbs.IDB
	seqValue         int32
	tmpNodes         *tmpNodeWatchers
	valueWatchers    *valueWatchers
	childrenWatchers *childrenWatchers
}

func NewDBR(c *xdb.DB, o *r.Options) (*DBR, error) {
	db, err := db.NewDB(c.Provider, c.ConnString, c.MaxOpen, c.MaxIdle, c.LifeTime)
	if err != nil {
		return nil, err
	}

	return &DBR{
		db:               db,
		seqValue:         10000,
		tmpNodes:         newTmpNodeWatchers(db),
		valueWatchers:    newValueWatchers(db),
		childrenWatchers: newChildrenWatchers(db),
	}, nil
}

//Close 关闭当前服务
func (r *DBR) Start() error {
	go r.valueWatchers.Start()
	go r.childrenWatchers.Start()
	go r.tmpNodes.Start()
	return nil
}

//Close 关闭当前服务
func (r *DBR) Close() error {
	r.valueWatchers.Close()
	r.childrenWatchers.Close()
	r.tmpNodes.Close()
	return nil
}

//dbrFactory 基于dbr的注册中心
type dbrFactory struct {
	proto string
	opts  *r.Options
}

//Create 根据配置生成dbr注册中心
func (z *dbrFactory) Create(opts ...r.Option) (r.IRegistry, error) {
	for i := range opts {
		opts[i](z.opts)
	}

	dbConf := &xdb.DB{
		Provider:   z.proto,
		ConnString: fmt.Sprintf("%s:%s@%s?charset=utf8", z.opts.Auth.Username, z.opts.Auth.Password, z.opts.Addrs[0]),
		MaxOpen:    10,
		MaxIdle:    10,
		LifeTime:   600,
	}

	r, err := NewDBR(dbConf, z.opts)
	if err != nil {
		return nil, err
	}
	r.Start()

	return r, err
}

var MYSQL = "mysql"
var ORACLE = "oracle"

func init() {
	r.Register(MYSQL, &dbrFactory{proto: MYSQL, opts: &r.Options{}})
	r.Register(ORACLE, &dbrFactory{proto: ORACLE, opts: &r.Options{}})

}
