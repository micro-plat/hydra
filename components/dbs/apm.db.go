package dbs

import (
	"fmt"
	"runtime"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/apm"
	"github.com/micro-plat/lib4go/db"
)

var _ IDB = &APMDB{}

type APMDB struct {
	orgdb    IDB
	name     string
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

func NewAPMDB(provider, dbName string, orgdb IDB) (IDB, error) {

	return &APMDB{
		name:     dbName,
		orgdb:    orgdb,
		provider: provider,
	}, nil
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

	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "db.Query", sqlkey, callback)
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
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "db.Scalar", sqlkey, callback)
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
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "db.Execute", sqlkey, callback)
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
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "db.Executes", sqlkey, callback)
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
	result := apmExecute(d.provider, d.name, "db.ExecuteSP", procName, callback)
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
	result := apmExecute(d.provider, d.name, "db.Begin", d.name, callback)
	return NewAPMDBTrans(d, result.DBTrans), result.Error
}
func (d *APMDB) Close() {
	d.orgdb.Close()
}

func (d *APMDB) GetProvider() string {
	return d.provider
}

func apmExecute(provider, name, operationName, sqlkey string, callback DBCallback) *CallResult {

	ctx := context.Current()
	apmCtx := ctx.APMContext()
	if apmCtx == nil {
		return callback()
	}
	rootCtx := apmCtx.GetRootCtx()
	tracer := apmCtx.GetTracer()

	span, err := tracer.CreateExitSpan(rootCtx, operationName, "db.Host", func(header string) error {
		return nil
	})
	if err != nil {
		err = fmt.Errorf("(DB/Trans).tracer.CreateExitSpan:%+v", err)
		return callback()
	}
	fmt.Println("apmExecute.3")
	defer span.End()
	//执行db 请求
	res := callback()
	span.SetComponent(apm.ComponentIDGODBClient)
	span.Tag("DBProvider", fmt.Sprintf("%s[%s]", provider, name))
	span.Tag("SQLKey", sqlkey)
	span.SetSpanLayer(apm.SpanLayer_Database)

	return res
}

func getSQLKey() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d=>%s", file, line, fn.Name())
}
