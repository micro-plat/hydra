package dbs

import (
	"fmt"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	"github.com/micro-plat/hydra/components/pkgs/apm/apmtypes"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/db"

	xdb "github.com/micro-plat/hydra/conf/vars/db"
)

var _ IDB = &APMDB{}

type APMDB struct {
	orgdb    IDB
	provider string
}

type DBCallback func() *CallResult
type CallResult struct {
	DBTrans  db.IDBTrans
	Error    error
	Args     []interface{}
	Query    string
	Rows     db.QueryRows
	EffCount int64
	LastID   int64
	Data     interface{}
}

func NewAPMDB(dbConf xdb.DB) (IDB, error) {

	orgdb, err := db.NewDB(dbConf.Provider,
		dbConf.ConnString,
		dbConf.MaxOpen,
		dbConf.MaxIdle,
		dbConf.LifeTime)
	if !global.Def.IsUseAPM() {
		return orgdb, err
	}

	return &APMDB{
		orgdb:    orgdb,
		provider: dbConf.Provider,
	}, err
}

func (d *APMDB) Query(sql string, input map[string]interface{}) (db.QueryRows, string, []interface{}, error) {

	callback := func() *CallResult {
		data, query, args, err := d.orgdb.Query(sql, input)
		return &CallResult{
			Rows:  data,
			Query: query,
			Args:  args,
			Error: err,
		}
	}
	result := apmExecute(d.provider, "db.Query", callback)
	return result.Rows, result.Query, result.Args, result.Error
}
func (d *APMDB) Scalar(sql string, input map[string]interface{}) (interface{}, string, []interface{}, error) {

	callback := func() *CallResult {
		data, query, args, err := d.orgdb.Scalar(sql, input)
		return &CallResult{
			Data:  data,
			Query: query,
			Args:  args,
			Error: err,
		}
	}
	result := apmExecute(d.provider, "db.Scalar", callback)
	return result.Data, result.Query, result.Args, result.Error
}
func (d *APMDB) Execute(sql string, input map[string]interface{}) (int64, string, []interface{}, error) {

	callback := func() *CallResult {
		effCount, query, args, err := d.orgdb.Execute(sql, input)
		return &CallResult{
			EffCount: effCount,
			Query:    query,
			Args:     args,
			Error:    err,
		}
	}
	result := apmExecute(d.provider, "db.Execute", callback)
	return result.EffCount, result.Query, result.Args, result.Error

}
func (d *APMDB) Executes(sql string, input map[string]interface{}) (int64, int64, string, []interface{}, error) {
	callback := func() *CallResult {
		lastID, effCount, query, args, err := d.orgdb.Executes(sql, input)
		return &CallResult{
			LastID:   lastID,
			EffCount: effCount,
			Query:    query,
			Args:     args,
			Error:    err,
		}
	}
	result := apmExecute(d.provider, "db.Executes", callback)
	return result.LastID, result.EffCount, result.Query, result.Args, result.Error

}
func (d *APMDB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (int64, string, error) {
	callback := func() *CallResult {
		effCount, query, err := d.orgdb.ExecuteSP(procName, input, output)
		return &CallResult{
			EffCount: effCount,
			Query:    query,
			Error:    err,
		}
	}
	result := apmExecute(d.provider, "db.ExecuteSP", callback)
	return result.EffCount, result.Query, result.Error

}

func (d *APMDB) Begin() (db.IDBTrans, error) {

	callback := func() *CallResult {
		dbTrans, err := d.orgdb.Begin()
		return &CallResult{
			DBTrans: dbTrans,
			Error:   err,
		}
	}
	result := apmExecute(d.provider, "db.Begin", callback)
	return NewAPMDBTrans(d, result.DBTrans), result.Error
}
func (d *APMDB) Close() {
	d.orgdb.Close()
}

func (d *APMDB) GetProvider() string {
	return d.provider
}

func apmExecute(provider, operationName string, callback DBCallback) *CallResult {
	ctx := context.Current()
	apmCfg := ctx.ServerConf().GetAPMConf()
	if apmCfg.Disable {
		return callback()
	}
	if !apmCfg.DB {
		return callback()
	}
	fmt.Println("apmExecute.1")
	tmp, ok := ctx.Meta().Get(apm.TraceInfo)
	if !ok {
		return callback()
	}
	fmt.Println("apmExecute.2")
	apmInfo := tmp.(*apm.APMInfo)
	rootCtx := apmInfo.RootCtx
	tracer := apmInfo.Tracer

	span, err := tracer.CreateExitSpan(rootCtx, operationName, "db.Host", func(header string) error {
		return nil
	})
	if err != nil {
		return callback()
	}
	fmt.Println("apmExecute.3")
	defer span.End()
	//执行db 请求
	res := callback()
	span.SetComponent(apmtypes.ComponentIDGODBClient)
	span.Tag("DBProvider", provider)
	span.SetSpanLayer(apm.SpanLayer_Database)

	return res
}
