package dbs

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/db"
)

var _ db.IDBTrans = &APMDBTrans{}

type APMDBTrans struct {
	orgTrans db.IDBTrans
	provider string
	name     string
}

func NewAPMDBTrans(dbObj *APMDB, dbTrans db.IDBTrans) db.IDBTrans {
	if !global.Def.IsUseAPM() {
		return dbTrans
	}
	return &APMDBTrans{
		orgTrans: dbTrans,
		name:     dbObj.name,
		provider: dbObj.provider,
	}
}

func (d *APMDBTrans) Query(sql string, input map[string]interface{}) (db.QueryRows, string, []interface{}, error) {

	callback := func() *CallResult {
		data, query, args, err := d.orgTrans.Query(sql, input)
		return &CallResult{
			Rows:  data,
			Query: query,
			Args:  args,
			Error: err,
		}
	}
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "Trans.Query", sqlkey, callback)
	return result.Rows, result.Query, result.Args, result.Error
}
func (d *APMDBTrans) Scalar(sql string, input map[string]interface{}) (interface{}, string, []interface{}, error) {

	callback := func() *CallResult {
		data, query, args, err := d.orgTrans.Scalar(sql, input)
		return &CallResult{
			Data:  data,
			Query: query,
			Args:  args,
			Error: err,
		}
	}
	sqlkey := getSQLKey()
	result := apmExecute(d.provider,  d.name,"Trans.Scalar", sqlkey, callback)
	return result.Data, result.Query, result.Args, result.Error
}
func (d *APMDBTrans) Execute(sql string, input map[string]interface{}) (int64, string, []interface{}, error) {

	callback := func() *CallResult {
		effCount, query, args, err := d.orgTrans.Execute(sql, input)
		return &CallResult{
			EffCount: effCount,
			Query:    query,
			Args:     args,
			Error:    err,
		}
	}
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "Trans.Execute", sqlkey, callback)
	return result.EffCount, result.Query, result.Args, result.Error

}
func (d *APMDBTrans) Executes(sql string, input map[string]interface{}) (int64, int64, string, []interface{}, error) {
	callback := func() *CallResult {
		lastID, effCount, query, args, err := d.orgTrans.Executes(sql, input)
		return &CallResult{
			LastID:   lastID,
			EffCount: effCount,
			Query:    query,
			Args:     args,
			Error:    err,
		}
	}
	sqlkey := getSQLKey()
	result := apmExecute(d.provider, d.name, "Trans.Executes", sqlkey, callback)
	return result.LastID, result.EffCount, result.Query, result.Args, result.Error

}
func (d *APMDBTrans) Rollback() error {

	callback := func() *CallResult {
		err := d.orgTrans.Rollback()
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.provider, d.name, "Trans.Rollback", d.name, callback)
	return result.Error

}
func (d *APMDBTrans) Commit() error {
	callback := func() *CallResult {
		err := d.orgTrans.Commit()
		return &CallResult{
			Error: err,
		}
	}
	result := apmExecute(d.provider, d.name, "Trans.Commit", d.name, callback)
	return result.Error
}

func (d *APMDBTrans) GetProvider() string {
	return d.provider
}
