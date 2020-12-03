package redis

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type value struct {
	IsTemp  bool   `json:"is_temp"`
	Data    []byte `json:"data"`
	Version int32  `json:"version"`
}

var s20200101, _ = time.Parse("20060102", "20200101")
var start int32 = 100

func newValue(data string, tmp bool) *value {
	return &value{
		Data:    []byte(data),
		IsTemp:  tmp,
		Version: int32(time.Now().Unix()-s20200101.Unix()) + atomic.AddInt32(&start, 1),
	}
}
func newValueByJSON(d string) (*value, error) {
	v := value{}
	err := json.Unmarshal([]byte(d), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *value) String() string {
	buf, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(buf)
}

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}
type childrenEntity struct {
	children []string
	version  int32
	path     string
	Err      error
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

func (v *childrenEntity) GetValue() ([]string, int32) {
	return v.children, v.version
}
func (v *childrenEntity) GetError() error {
	return v.Err
}
func (v *childrenEntity) GetPath() string {
	return v.path
}
