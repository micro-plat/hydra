package conf

import (
	"sync"
)

type metadata struct {
	data map[string]interface{}
	lock sync.RWMutex
}

func (m *metadata) Get(key string) (interface{}, bool) {
	m.lock.RLocker().Lock()
	defer m.lock.RLocker().Unlock()

	data, ok := m.data[key]
	return data, ok
}
func (m *metadata) Set(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[key] = value
}
func (m *metadata) CSet(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = value
}

type IMetadata interface {
	Get(key string) (interface{}, bool)
	Set(key string, v interface{})
}

type Metadata struct {
	Name string
	Type string
	data *metadata
}

func NewMetadata(name, tp string) *Metadata {
	return &Metadata{
		Name: name,
		Type: tp,
		data: &metadata{
			data: make(map[string]interface{}),
		},
	}
}

func (s *Metadata) Get(key string) (interface{}, bool) {
	return s.data.Get(key)
}
func (s *Metadata) Set(key string, v interface{}) {
	s.data.Set(key, v)
}

func (s *Metadata) CSet(key string, v interface{}) {
	s.data.CSet(key, v)
}
