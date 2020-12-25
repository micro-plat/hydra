package global

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//db 数据库处理逻辑
type db struct {
	sqls     []string
	handlers []func() error
}

//AddBSQL 添加执行SQL
func (d *db) AddBSQL(sqls ...[]byte) {
	for _, sql := range sqls {
		nsql := strings.Split(strings.Trim(types.BytesToString(sql), ";"), ";")
		for _, m := range nsql {
			if strings.TrimSpace(m) != "" {
				d.sqls = append(d.sqls, strings.TrimSpace(m))
			}
		}
	}
}

//AddBSQL 添加执行SQL
func (d *db) AddSQL(sqls ...string) {
	for _, sql := range sqls {
		nsql := strings.Split(strings.Trim(sql, ";"), ";")
		for _, m := range nsql {
			if strings.TrimSpace(m) != "" {
				d.sqls = append(d.sqls, strings.TrimSpace(m))
			}
		}

	}

}

//AddHandler 添加处理函数
func (d *db) AddHandler(fs ...interface{}) {
	for _, fn := range fs {
		var nfunc func() error
		hasMatch := false
		if fx, ok := fn.(func()); ok {
			hasMatch = true
			nfunc = func() error {
				fx()
				return nil
			}
		}
		if fx, ok := fn.(func() error); ok {
			hasMatch = true
			nfunc = fx
		}
		if !hasMatch {
			panic("函数签名格式不正确，支持的格式有func(){} 或 func()error{}")
		}
		d.handlers = append(d.handlers, nfunc)
	}
}

//GetSQLs 获取所有SQL语句
func (d *db) GetSQLs() []string {
	return d.sqls
}

//GetHandlers 获取所有处理函数
func (d *db) GetHandlers() []func() error {
	return d.handlers
}
