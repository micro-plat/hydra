package lnfs

import (
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//eDirEntity 目录实体
type eDirEntity struct {
	Path    string    `json:"path,omitempty" valid:"required"`
	Name    string    `json:"name,omitempty" valid:"required"`
	ModTime time.Time `json:"modTime,omitempty" valid:"required"`
	Size    int64     `json:"buffer,omitempty" valid:"required"`
}

//目录实体列表
type eDirEntityList []*eDirEntity

func (m eDirEntityList) GetMap() map[string]interface{} {
	x := make(map[string]interface{})
	for _, v := range m {
		x[v.Path] = v
	}
	return x
}

func (m eDirEntityList) Equal(c cmap.ConcurrentMap) bool {
	if len(m) != c.Count() {
		return false
	}
	vmap := m.GetMap()
	for k := range vmap {
		if !c.Has(k) {
			return false
		}
	}
	return true
}
