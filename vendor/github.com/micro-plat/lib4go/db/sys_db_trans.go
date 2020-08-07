package db

import "database/sql"

//SysDBTransaction 事务
type SysDBTransaction struct {
	tx *sql.Tx
}

//Query 执行查询
func (t *SysDBTransaction) Query(query string, args ...interface{}) (dataRows QueryRows, colus []string, err error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	dataRows, colus, err = resolveRows(rows, 0)
	return
}

//Executes 执行SQL操作语句
func (t *SysDBTransaction) Executes(query string, args ...interface{}) (lastInsertId, affectedRow int64, err error) {
	result, err := t.tx.Exec(query, args...)
	if err != nil {
		return
	}
	lastInsertId, err = result.LastInsertId()
	affectedRow, err = result.RowsAffected()
	return
}

//Execute 执行SQL操作语句
func (t *SysDBTransaction) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := t.tx.Exec(query, args...)
	if err != nil {
		return
	}
	affectedRow, err = result.RowsAffected()
	return
}

//Rollback 回滚所有操作
func (t *SysDBTransaction) Rollback() error {
	return t.tx.Rollback()
}

//Commit 提交所有操作
func (t *SysDBTransaction) Commit() error {
	return t.tx.Commit()
}
