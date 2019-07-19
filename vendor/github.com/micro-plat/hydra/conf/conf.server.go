package conf

import (
	"sync"
)

type metadata struct {
	data map[string]interface{}
	lock sync.RWMutex
}

func (m *metadata) Get(key string) interface{} {
	m.lock.RLocker().Lock()
	defer m.lock.RLocker().Unlock()

	data := m.data[key]
	return data
}
func (m *metadata) Set(key string, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = value
}

type MetadataConf struct {
	Name     string
	Type     string
	metadata metadata
}

func (s *MetadataConf) GetMetadata(key string) interface{} {
	return s.metadata.Get(key)
}
func (s *MetadataConf) SetMetadata(key string, v interface{}) {
	s.metadata.Set(key, v)
}
