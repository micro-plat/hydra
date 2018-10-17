package db

import (
	"time"

	"github.com/micro-plat/lib4go/types"
)

type IQueryRow interface {
	GetString(name string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	Has(name string) bool
	GetMustString(name string) (string, bool)
	GetMustInt(name string) (int, bool)
	GetMustFloat32(name string) (float32, bool)
	GetMustFloat64(name string) (float64, bool)
	GetDatetime(name string, format ...string) (time.Time, error)
	ToStruct(o interface{}) error
}

//QueryRow 查询的数据行
type QueryRow = types.XMap

//QueryRows 多行数据
type QueryRows []QueryRow

//ToStruct 将当前对象转换为指定的struct
func (q QueryRows) ToStruct(o interface{}) error {
	return types.Map2Struct(q, o)
}

//IsEmpty 当前数据集是否为空
func (q QueryRows) IsEmpty() bool {
	return q == nil || len(q) == 0
}

//Len 获取当前数据集的长度
func (q QueryRows) Len() int {
	return len(q)
}

//Get 获取指定索引的数据
func (q QueryRows) Get(i int) QueryRow {
	if q == nil || i >= len(q) || i < 0 {
		return QueryRow{}
	}
	return q[i]
}
