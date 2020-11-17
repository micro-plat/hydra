package tgo

import (
	"github.com/d5/tengo/v2"
	"github.com/micro-plat/lib4go/types"
)

func var2Mpa(s *tengo.Compiled) types.XMap {
	mp := types.NewXMap()
	vb := s.GetAll()
	for _, v := range vb {
		mp[v.Name()] = v.Value()
	}
	return mp
}
