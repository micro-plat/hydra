package db

import (
	"github.com/micro-plat/lib4go/types"
)

//QueryRow 单行数据
type QueryRow = types.XMap

//NewQueryRow 构建QueryRow对象
func NewQueryRow(len ...int) QueryRow {
	return types.NewXMap(len...)
}

//QueryRows 多行数据
type QueryRows = types.XMaps

//NewQueryRows 构建QueryRows
func NewQueryRows(len ...int) QueryRows {
	return types.NewXMaps(len...)
}
