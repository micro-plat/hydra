package context

import "strings"

type meta struct {
	d map[string]interface{}
}

func (m *meta) Keys() []string {
	keys := make([]string, 0, len(m.d))
	for k := range m.d {
		keys = append(keys, k)
	}
	return keys
}

func (m *meta) Get(name string) (interface{}, bool) {
	v, ok := m.d[name]
	return v, ok
}

//Meta 元数据
type Meta struct {
	inputParams
	m *meta
}

//NewMeta 创建meta
func NewMeta() *Meta {
	m := &meta{d: make(map[string]interface{})}
	return &Meta{
		m: m,
		inputParams: inputParams{
			data: m,
		},
	}
}

//Sets 设置map参数
func (m *Meta) Sets(i map[string]interface{}) {
	m.m.d = i
}

//SetStrings 设置map参数
func (m *Meta) SetStrings(input map[string]string) {
	for i, v := range input {
		m.m.d[strings.ToLower(i)] = v
	}
}

//Set 设置元数据
func (m *Meta) Set(key string, value interface{}) {
	m.m.d[key] = value
}

//Get 获取元数据
func (m *Meta) Get(name string) (interface{}, bool) {
	c, ok := m.data.Get(name)
	return c, ok
}
