package db

import (
	"database/sql"
	"time"

	"github.com/micro-plat/lib4go/db/tpl"
)

// IDB 数据库操作接口,安装可需能需要执行export LD_LIBRARY_PATH=/usr/local/lib
type IDB interface {
	IDBExecuter
	ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (row int64, err error)
	Begin() (IDBTrans, error)
	Close()
}

// IDBTrans 数据库事务接口
type IDBTrans interface {
	IDBExecuter
	Rollback() error
	Commit() error
}

// IDBExecuter 数据库操作对象集合
type IDBExecuter interface {
	FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error)
	Query(sql string, input map[string]interface{}) (data QueryRows, err error)
	QueryBatch(sql []string, input map[string]interface{}) (data QueryRows, err error)
	Scalar(sql string, input map[string]interface{}) (data interface{}, err error)
	Execute(sql string, input map[string]interface{}) (row int64, err error)
	Executes(sql string, input map[string]interface{}) (lastInsertID int64, affectedRow int64, err error)
	ExecuteBatch(sql []string, input map[string]interface{}) (data QueryRows, err error)
	InsertBatch(sql string, inputs []map[string]interface{}) (row int64, err error)
	UpdateBatch(sql string, inputs []map[string]interface{}) (row int64, err error)
}

// DB 数据库操作类
type DB struct {
	db  ISysDB
	tpl tpl.ITPLContext
}

// NewDB 创建DB实例
func NewDB(provider string, connString string, maxOpen int, maxIdle int, maxLifeTime int) (obj *DB, err error) {
	obj = &DB{}
	if obj.tpl, err = tpl.GetDBContext(provider); err != nil {
		return
	}
	obj.db, err = NewSysDB(provider, connString, maxOpen, maxIdle, time.Duration(maxLifeTime)*time.Second)
	return
}

// GetTPL 获取模板翻译参数
func (db *DB) GetTPL() tpl.ITPLContext {
	return db.tpl
}
func (db *DB) FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error) {
	return fetchRows(db.db, db.tpl, sql, input)
}

// Query 查询数据
func (db *DB) Query(sql string, input map[string]interface{}) (data QueryRows, err error) {
	return query(db.db, db.tpl, sql, input)
}

// Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DB) Scalar(sql string, input map[string]interface{}) (data interface{}, err error) {
	return scalar(db.db, db.tpl, sql, input)
}

// Executes 根据包含@名称占位符的语句执行查询语句
func (db *DB) Executes(sql string, input map[string]interface{}) (insertID int64, row int64, err error) {
	return executes(db.db, db.tpl, sql, input)
}

// Execute 根据包含@名称占位符的语句执行查询语句
func (db *DB) Execute(sql string, input map[string]interface{}) (row int64, err error) {
	return execute(db.db, db.tpl, sql, input)
}

func (db *DB) InsertBatch(sql string, inputs []map[string]interface{}) (row int64, err error) {
	return insertBatch(db.db, db.tpl, sql, inputs)
}
func (db *DB) UpdateBatch(sql string, inputs []map[string]interface{}) (row int64, err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}
		tx.Rollback()
	}()

	return updateSave(tx, db.tpl, sql, inputs)
}

// ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (row int64, err error) {
	return executeSP(db.db, db.tpl, procName, input, output...)
}

// ExecuteBatch 批量执行SQL语句
func (db *DB) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(db, sqls, input)
}

// QueryBatch 批量执行多个查询语句
func (db *DB) QueryBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return queryBatch(db, sqls, input)
}

// Replace 替换SQL语句中的参数
func (db *DB) Replace(sql string, args []interface{}) string {
	return db.tpl.Replace(sql, args)
}

// Begin 创建事务
func (db *DB) Begin() (t IDBTrans, err error) {
	tt := &DBTrans{tpl: db.tpl}
	if tt.tx, err = db.db.Begin(); err != nil {
		return
	}
	return tt, nil
}

// Close  关闭当前数据库连接
func (db *DB) Close() {
	db.db.Close()
}
