package conf

import (
	"time"

	"github.com/micro-plat/lib4go/types"
)

var _ IMeta = meta{}

type IMeta interface {
	Keys() []string
	Get(name string) (interface{}, bool)
	GetValue(name string) interface{}
	GetString(name string) string
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	GetMustString(name string) (string, bool)
	GetMustInt(name string) (int, bool)
	GetMustFloat32(name string) (float32, bool)
	GetMustFloat64(name string) (float64, bool)
	GetDatetime(name string, format ...string) (time.Time, error)

	Set(name string, value interface{})
	Has(name string) bool
	IsEmpty() bool
	Len() int
	ToStruct(o interface{}) error
}

type meta types.XMap

func (m meta) Set(name string, value interface{}) {
	v := types.IXMap(meta)
	v.SetValue(name, value)
}
