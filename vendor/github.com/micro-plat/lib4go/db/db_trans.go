package db

import (
	"github.com/micro-plat/lib4go/db/tpl"
)

//DBTrans 数据库事务操作类
type DBTrans struct {
	tpl tpl.ITPLContext
	tx  ISysDBTrans
}

//Query 查询数据
func (t *DBTrans) Query(sql string, input map[string]interface{}) (data QueryRows, err error) {
	query, args := t.tpl.GetSQLContext(sql, input)
	data, err = t.tx.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (t *DBTrans) Scalar(sql string, input map[string]interface{}) (data interface{}, err error) {
	query, args := t.tpl.GetSQLContext(sql, input)
	result, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	if result.Len() == 0 || result.Get(0).IsEmpty() {
		return nil, nil
	}
	data, _ = result.Get(0).Get(result.Get(0).Keys()[0])
	return
}

//Executes 执行SQL操作语句
func (t *DBTrans) Executes(sql string, input map[string]interface{}) (lastInsertID, affectedRow int64, err error) {
	query, args := t.tpl.GetSQLContext(sql, input)
	lastInsertID, affectedRow, err = t.tx.Executes(query, args...)
	if err != nil {
		return 0, 0, getDBError(err, query, args)
	}
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) Execute(sql string, input map[string]interface{}) (row int64, err error) {
	query, args := t.tpl.GetSQLContext(sql, input)
	row, err = t.tx.Execute(query, args...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) ExecuteSP(sql string, input map[string]interface{}) (row int64, err error) {
	query, args := t.tpl.GetSPContext(sql, input)
	row, err = t.tx.Execute(query, args...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}

//ExecuteBatch 批量执行SQL语句
func (t *DBTrans) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(t, sqls, input)
}

//Rollback 回滚所有操作
func (t *DBTrans) Rollback() error {
	return t.tx.Rollback()
}

//Commit 提交所有操作
func (t *DBTrans) Commit() error {
	return t.tx.Commit()
}
