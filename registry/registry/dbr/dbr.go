package dbr

import (
	"github.com/micro-plat/hydra/components/dbs"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/db"
)

type DBR struct {
	db               dbs.IDB
	sqltexture       *sqltexture
	seqValue         int32
	tmpNodes         *tmpNodeWatchers
	valueWatchers    *valueWatchers
	childrenWatchers *childrenWatchers
}

func NewDBR(c *xdb.DB, sqltexture *sqltexture, o *r.Options) (*DBR, error) {
	db, err := db.NewDB(c.Provider, c.ConnString, c.MaxOpen, c.MaxIdle, c.LifeTime)
	if err != nil {
		return nil, err
	}

	return &DBR{
		db:               db,
		seqValue:         10000,
		sqltexture:       sqltexture,
		tmpNodes:         newTmpNodeWatchers(db, sqltexture),
		valueWatchers:    newValueWatchers(db, sqltexture),
		childrenWatchers: newChildrenWatchers(db, sqltexture),
	}, nil
}

//Close 关闭当前服务
func (r *DBR) Start() error {
	if err := r.CreateStructure(); err != nil {
		return err
	}
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

	var dbConf *xdb.DB
	var sqltexture *sqltexture
	switch z.proto {
	case MYSQL:
		sqltexture = &mysqltexture
		dbConf = mysql.NewBy(z.opts.Auth.Username, z.opts.Auth.Password, z.opts.Addrs[0], z.opts.Metadata["db"], xdb.WithConnect(10, 6, 900))
	case ORACLE:
		dbConf = oracle.NewBy(z.opts.Auth.Username, z.opts.Auth.Password, z.opts.Addrs[0], xdb.WithConnect(10, 6, 900))

	}

	r, err := NewDBR(dbConf, sqltexture, z.opts)
	if err != nil {
		return nil, err
	}
	if err := r.Start(); err != nil {
		return nil, err
	}

	return r, err
}

var MYSQL = "mysql"
var ORACLE = "oracle"

func init() {
	r.Register(MYSQL, &dbrFactory{proto: MYSQL, opts: &r.Options{}})
	// r.Register(ORACLE, &dbrFactory{proto: ORACLE, opts: &r.Options{}})//暂未提供SQL

}
