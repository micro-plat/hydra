package db

import (
	"time"

	"github.com/micro-plat/lib4go/db/tpl"
)

//IDB 数据库操作接口,安装可需能需要执行export LD_LIBRARY_PATH=/usr/local/lib
type IDB interface {
	IDBExecuter
	ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (row int64, err error)
	Begin() (IDBTrans, error)
	Close()
}

//IDBTrans 数据库事务接口
type IDBTrans interface {
	IDBExecuter
	Rollback() error
	Commit() error
}

//IDBExecuter 数据库操作对象集合
type IDBExecuter interface {
	Query(sql string, input map[string]interface{}) (data QueryRows, err error)
	Scalar(sql string, input map[string]interface{}) (data interface{}, err error)
	Execute(sql string, input map[string]interface{}) (row int64, err error)
	Executes(sql string, input map[string]interface{}) (lastInsertID int64, affectedRow int64, err error)
	ExecuteBatch(sql []string, input map[string]interface{}) (data QueryRows, err error)
}

//DB 数据库操作类
type DB struct {
	db  ISysDB
	tpl tpl.ITPLContext
}

//NewDB 创建DB实例
func NewDB(provider string, connString string, maxOpen int, maxIdle int, maxLifeTime int) (obj *DB, err error) {
	obj = &DB{}
	obj.tpl, err = tpl.GetDBContext(provider)
	if err != nil {
		return
	}
	obj.db, err = NewSysDB(provider, connString, maxOpen, maxIdle, time.Duration(maxLifeTime)*time.Second)
	return
}

//GetTPL 获取模板翻译参数
func (db *DB) GetTPL() tpl.ITPLContext {
	return db.tpl
}

//Query 查询数据
func (db *DB) Query(sql string, input map[string]interface{}) (data QueryRows, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	data, err = db.db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DB) Scalar(sql string, input map[string]interface{}) (data interface{}, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	result, err := db.db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	if result.Len() == 0 || result.Get(0).IsEmpty() {
		return nil, nil
	}
	data, _ = result.Get(0).Get(result.Get(0).Keys()[0])
	return
}

//Executes 根据包含@名称占位符的语句执行查询语句
func (db *DB) Executes(sql string, input map[string]interface{}) (insertID int64, row int64, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	insertID, row, err = db.db.Executes(query, args...)
	if err != nil {
		return 0, 0, getDBError(err, query, args)
	}
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (db *DB) Execute(sql string, input map[string]interface{}) (row int64, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	row, err = db.db.Execute(query, args...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (row int64, err error) {
	query, args := db.tpl.GetSPContext(procName, input)
	ni := append(args, output...)
	row, err = db.db.Execute(query, ni...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}

//ExecuteBatch 批量执行SQL语句
func (db *DB) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(db, sqls, input)
}

//Replace 替换SQL语句中的参数
func (db *DB) Replace(sql string, args []interface{}) string {
	return db.tpl.Replace(sql, args)
}

//Begin 创建事务
func (db *DB) Begin() (t IDBTrans, err error) {
	tt := &DBTrans{}
	tt.tx, err = db.db.Begin()
	if err != nil {
		return
	}
	tt.tpl = db.tpl
	return tt, nil
}

//Close  关闭当前数据库连接
func (db *DB) Close() {
	db.db.Close()
}
