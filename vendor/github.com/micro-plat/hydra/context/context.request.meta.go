package context

type meta struct {
	d map[string]interface{}
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

//Set 设置元数据
func (m *Meta) Set(key string, value interface{}) {
	m.m.d[key] = value
}

//Get 获取元数据
func (m *Meta) Get(name string) (interface{}, bool) {
	c, ok := m.data.Get(name)
	return c, ok
}
