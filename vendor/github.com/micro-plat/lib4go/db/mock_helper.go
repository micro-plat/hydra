// Package db 提供 Mock 辅助工具
package db

import (
	"github.com/micro-plat/lib4go/types"
)

// NewMockQueryRows 创建 Mock 查询结果
// 从 map 切片创建 QueryRows
func NewMockQueryRows(data []map[string]interface{}) QueryRows {
	if data == nil {
		return make(types.XMaps, 0)
	}
	rows := make(types.XMaps, 0, len(data))
	for _, d := range data {
		rows = append(rows, types.XMap(d))
	}
	return rows
}

// NewMockQueryRow 创建 Mock 单行结果
// 从 map 创建 QueryRow
func NewMockQueryRow(data map[string]interface{}) QueryRow {
	if data == nil {
		return types.XMap{}
	}
	return types.XMap(data)
}
