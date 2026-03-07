package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/types"
	//_ "github.com/mattn/go-oci8"
	//_ "github.com/mattn/go-sqlite3"
	//_ "gopkg.in/rana/ora.v4"
)

/*
github.com/mattn/go-oci8

http://www.simonzhang.net/?p=2890
http://blog.sina.com.cn/s/blog_48c95a190102w2ln.html
http://www.tudou.com/programs/view/yet9OngrV_4/
https://github.com/wendal/go-oci8/downloads
https://github.com/wendal/go-oci8

安装方法
1. 下载：http://www.oracle.com/technetwork/database/features/instant-client/index.html
2. 解压文件 unzip instantclient-basic-linux.x64-12.1.0.1.0.zip -d /usr/local/
3. 配置环境变量
vi .bash_profile
export ora_home=/usr/local/instantclient_12_1
export PATH=$PATH:$ora_home
export LD_LIBRARY_PATH=$ora_home


*/

const (
	//SQLITE3 Sqlite3数据库
	SQLITE3 = "sqlite3"
	//OCI8 oralce数据库
	OCI8 = "oci8"
	//ORA  oralce数据库
	ORA = "ora"
)

type IBaseDB interface {
	FetchRows(string, ...interface{}) (*sql.Rows, error)
	Query(string, ...interface{}) (QueryRows, error)
	Execute(string, ...interface{}) (int64, error)
	Executes(string, ...interface{}) (int64, int64, error)
}

type ISysDB interface {
	FetchRows(string, ...interface{}) (*sql.Rows, error)
	Query(string, ...interface{}) (QueryRows, error)
	Execute(string, ...interface{}) (int64, error)
	Executes(string, ...interface{}) (int64, int64, error)
	Begin() (ISysDBTrans, error)
	Close()
}

// ISysDBTrans 数据库事务接口
type ISysDBTrans interface {
	FetchRows(string, ...interface{}) (*sql.Rows, error)
	Query(string, ...interface{}) (QueryRows, error)
	Execute(string, ...interface{}) (int64, error)
	Executes(query string, args ...interface{}) (lastInsertID, affectedRow int64, err error)
	Rollback() error
	Commit() error
}

// SysDB 数据库实体
type SysDB struct {
	provider   string
	connString string
	db         *sql.DB
}

// NewSysDB 创建DB实例
func NewSysDB(provider string, connString string, maxOpen int, maxIdle int, maxLifeTime time.Duration) (obj *SysDB, err error) {
	if provider == "" || connString == "" {
		err = errors.New("provider 和 connString 不能为空")
		return
	}
	obj = &SysDB{provider: provider, connString: connString}
	switch strings.ToLower(provider) {
	case "ora", "oracle":
		obj.db, err = sql.Open(OCI8, connString)
	case "sqlite":
		obj.db, err = sql.Open(SQLITE3, connString)
	default:
		obj.db, err = sql.Open(provider, connString)
	}
	if err != nil {
		return
	}
	obj.db.SetMaxIdleConns(maxIdle)
	obj.db.SetMaxOpenConns(maxOpen)
	obj.db.SetConnMaxLifetime(maxLifeTime)
	err = obj.db.Ping()
	return
}
func (db *SysDB) FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// Query 执行SQL查询语句
func (db *SysDB) Query(query string, args ...interface{}) (dataRows QueryRows, err error) {
	rows, err := db.db.Query(query, args...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return
	}
	defer rows.Close()
	dataRows, err = resolveFullRows(rows)
	return

}

func resolveRows(rows *sql.Rows, col int) (dataRows QueryRows, columns []string, err error) {
	dataRows = NewQueryRows()
	colus, err := rows.Columns()
	if err != nil {
		return
	}
	columns = make([]string, 0, len(colus))
	for _, v := range colus {
		columns = append(columns, strings.ToLower(v))
	}

	for rows.Next() {
		row := types.NewXMap(len(columns))
		dataRows.Append(row)

		// 使用更通用的类型来扫描数据，以处理NULL值
		buffer := make([]interface{}, len(columns))
		values := make([]sql.NullString, len(columns))
		for i := range buffer {
			buffer[i] = &values[i]
		}

		if err = rows.Scan(buffer...); err != nil {
			return
		}

		for index := 0; index < len(columns) && (index < col || col == 0); index++ {
			key := columns[index]
			if values[index].Valid {
				row.SetValue(key, values[index].String)
			} else {
				row.SetValue(key, "") // 将NULL值设置为空字符串
			}
		}
	}
	return
}

// func resolveRows(rows *sql.Rows, col int) (dataRows QueryRows, columns []string, err error) {
// 	fmt.Println("-----------------------xxx-------------------------")
// 	dataRows = NewQueryRows()
// 	colus, err := rows.Columns()
// 	if err != nil {
// 		return
// 	}
// 	columns = make([]string, 0, len(colus))
// 	for _, v := range colus {
// 		columns = append(columns, strings.ToLower(v))
// 	}

// 	for rows.Next() {
// 		row := types.NewXMap(len(columns))
// 		dataRows.Append(row)
// 		var buffer []interface{}
// 		for index := 0; index < len(columns); index++ {
// 			var va []byte
// 			buffer = append(buffer, &va)
// 		}
// 		err = rows.Scan(buffer...)
// 		if err != nil {
// 			return
// 		}
// 		for index := 0; index < len(columns) && (index < col || col == 0); index++ {
// 			key := columns[index]
// 			value := buffer[index]
// 			if value == nil {
// 				continue
// 			} else {
// 				fmt.Println("0------------------------------------------")
// 				buff := value.(*[]byte)
// 				row[key] = bytes.NewBuffer(*buff).String()
// 				row[key] = strings.TrimPrefix(fmt.Sprintf("%s", value), "&")

// 				// row.SetValue(key, strings.TrimPrefix(fmt.Sprintf("%s", value), "&"))
// 			}
// 		}
// 	}
// 	return
// }

// Executes 执行SQL操作语句
func (db *SysDB) Executes(query string, args ...interface{}) (lastInsertID, affectedRow int64, err error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return
	}
	affectedRow, err = result.RowsAffected()
	return
}

// Execute 执行SQL操作语句
func (db *SysDB) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

// Begin 创建一个事务请求
func (db *SysDB) Begin() (r ISysDBTrans, err error) {
	t := &SysDBTransaction{}
	t.tx, err = db.db.Begin()
	return t, err
}

// Close 关闭数据库
func (db *SysDB) Close() {
	db.db.Close()
}
