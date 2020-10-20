package localmemory

import (
	"sync/atomic"
	"time"
)

type value struct {
	data    string
	version int32
}

var s20200101, _ = time.Parse("20060102", "20200101")
var start int32 = 100

func newValue(data string) *value {
	return &value{
		data:    data,
		version: int32(time.Now().Unix()-s20200101.Unix()) + atomic.AddInt32(&start, 1),
	}
}

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}
type valuesEntity struct {
	values  []string
	version int32
	path    string
	Err     error
}

func (v *valueEntity) GetPath() string {
	return v.path
}
func (v *valueEntity) GetValue() ([]byte, int32) {
	return v.Value, v.version
}
func (v *valueEntity) GetError() error {
	return v.Err
}

func (v *valuesEntity) GetValue() ([]string, int32) {
	return v.values, v.version
}
func (v *valuesEntity) GetError() error {
	return v.Err
}
func (v *valuesEntity) GetPath() string {
	return v.path
}
