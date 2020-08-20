package apm

import (
	//"fmt"

	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/pkgs/apm"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

const (
	//APMTypeNode 缓存在var配置中的类型名称
	APMTypeNode = "apm"

	//APMNameNode 缓存名称在var配置中的末节点名称
	APMNameNode = "apm"
)

//StandardAPM apm
type StandardAPM struct {
	c container.IContainer
}

//NewStandardAPM 创建APM
func NewStandardAPM(c container.IContainer) IComponentAPM {
	return &StandardAPM{c: c}
}

//GetRegularAPM 获取正式的没有异常缓存实例
func (s *StandardAPM) GetRegularAPM(instance string, names ...string) (c IAPM) {
	c, err := s.GetAPM(instance, names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetAPM 获取缓存操作对象
func (s *StandardAPM) GetAPM(instance string, names ...string) (c IAPM, err error) {
	//fmt.Println("instance:", instance)
	name := types.GetStringByIndex(names, 0, APMNameNode)
	obj, err := s.c.GetOrCreate(APMTypeNode, name, func(js *conf.JSONConf) (interface{}, error) {
		//fmt.Println("JSONConf:%s", string(js.GetRaw()))
		return apm.New(js.GetString("apmtype"), instance, string(js.GetRaw()))
	})
	if err != nil {
		return nil, err
	}
	return obj.(IAPM), nil
}
