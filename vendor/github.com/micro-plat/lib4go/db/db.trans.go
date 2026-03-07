package db

import (
	"database/sql"

	"github.com/micro-plat/lib4go/db/tpl"
)

// DBTrans 数据库事务操作类
type DBTrans struct {
	tpl tpl.ITPLContext
	tx  ISysDBTrans
}

func (t *DBTrans) FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error) {
	return fetchRows(t.tx, t.tpl, sql, input)
}

// Query 查询数据
func (t *DBTrans) Query(sql string, input map[string]interface{}) (data QueryRows, err error) {
	return query(t.tx, t.tpl, sql, input)
}

// Scalar 根据包含@名称占位符的查询语句执行查询语句
func (t *DBTrans) Scalar(sql string, input map[string]interface{}) (data interface{}, err error) {
	return scalar(t.tx, t.tpl, sql, input)
}

// Executes 执行SQL操作语句
func (t *DBTrans) Executes(sql string, input map[string]interface{}) (lastInsertID, affectedRow int64, err error) {
	return executes(t.tx, t.tpl, sql, input)
}

// Execute 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) Execute(sql string, input map[string]interface{}) (row int64, err error) {
	return execute(t.tx, t.tpl, sql, input)
}

func (db *DBTrans) InsertBatch(sql string, inputs []map[string]interface{}) (row int64, err error) {
	return insertBatch(db.tx, db.tpl, sql, inputs)
}
func (db *DBTrans) UpdateBatch(sql string, inputs []map[string]interface{}) (row int64, err error) {
	return updateSave(db.tx, db.tpl, sql, inputs)
}

// ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) ExecuteSP(sql string, input map[string]interface{}) (row int64, err error) {
	return executeSP(t.tx, t.tpl, sql, input)
}

// ExecuteBatch 批量执行SQL语句
func (t *DBTrans) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(t, sqls, input)
}

// ExecuteBatch 批量执行SQL语句
func (t *DBTrans) QueryBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return queryBatch(t, sqls, input)
}

// Rollback 回滚所有操作
func (t *DBTrans) Rollback() error {
	return t.tx.Rollback()
}

// Commit 提交所有操作
func (t *DBTrans) Commit() error {
	return t.tx.Commit()
}
