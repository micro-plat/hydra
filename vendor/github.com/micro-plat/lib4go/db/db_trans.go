package db

import "github.com/micro-plat/lib4go/db/tpl"

//DBTrans 数据库事务操作类
type DBTrans struct {
	tpl tpl.ITPLContext
	tx  ISysDBTrans
}

//Query 查询数据
func (t *DBTrans) Query(sql string, input map[string]interface{}) (data QueryRows, query string, args []interface{}, err error) {
	query, args = t.tpl.GetSQLContext(sql, input)
	data, _, err = t.tx.Query(query, args...)
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (t *DBTrans) Scalar(sql string, input map[string]interface{}) (data interface{}, query string, args []interface{}, err error) {
	query, args = t.tpl.GetSQLContext(sql, input)
	result, colus, err := t.tx.Query(query, args...)
	if err != nil || len(result) == 0 || len(result[0]) == 0 || len(colus) == 0 {
		return
	}
	data = result[0][colus[0]]
	return
}

//Executes 执行SQL操作语句
func (t *DBTrans) Executes(sql string, input map[string]interface{}) (lastInsertId, affectedRow int64, query string, args []interface{}, err error) {
	query, args = t.tpl.GetSQLContext(sql, input)
	lastInsertId, affectedRow, err = t.tx.Executes(query, args...)
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) Execute(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = t.tpl.GetSQLContext(sql, input)
	row, err = t.tx.Execute(query, args...)
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (t *DBTrans) ExecuteSP(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = t.tpl.GetSPContext(sql, input)
	row, err = t.tx.Execute(query, args...)
	return
}

//Rollback 回滚所有操作
func (t *DBTrans) Rollback() error {

	return t.tx.Rollback()
}

//Commit 提交所有操作
func (t *DBTrans) Commit() error {
	return t.tx.Commit()
}
