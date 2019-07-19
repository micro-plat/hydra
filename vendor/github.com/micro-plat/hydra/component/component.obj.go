package component

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var _ IComponentGlobalVarObject = &GlobalVarObjectCache{}

//IComponentGlobalVarObject Component Cache
type IComponentGlobalVarObject interface {
	GetGlobalObject(tpName string, name string) (c interface{}, err error)
	SaveGlobalObject(tpName string, name string, f func(c conf.IConf) (interface{}, error)) (bool, interface{}, error)
	Close() error
}

//GlobalVarObjectCache cache
type GlobalVarObjectCache struct {
	IContainer
	cacheMap  cmap.ConcurrentMap
	closeList []CloseHandler
}

//NewGlobalVarObjectCache 创建cache
func NewGlobalVarObjectCache(c IContainer) *GlobalVarObjectCache {
	return &GlobalVarObjectCache{IContainer: c, cacheMap: cmap.New(2), closeList: make([]CloseHandler, 0, 1)}
}

//GetGlobalObject 获取全局对象
func (s *GlobalVarObjectCache) GetGlobalObject(tpName string, name string) (c interface{}, err error) {
	cacheConf, err := s.IContainer.GetVarConf(tpName, name)
	if err != nil {
		return nil, fmt.Errorf("%s %v", registry.Join("/", s.GetPlatName(), "var", tpName, name), err)
	}
	key := fmt.Sprintf("%s/%s:%d", tpName, name, cacheConf.GetVersion())
	c, ok := s.cacheMap.Get(key)
	if !ok {
		err = fmt.Errorf("缓存对象未创建:%s", registry.Join("/", s.GetPlatName(), "var", tpName, name))
		return
	}
	return c, nil
}

//SaveGlobalObject 缓存全局对象
func (s *GlobalVarObjectCache) SaveGlobalObject(tpName string, name string, f func(c conf.IConf) (interface{}, error)) (bool, interface{}, error) {
	cacheConf, err := s.IContainer.GetVarConf(tpName, name)
	if err != nil {
		return false, nil, fmt.Errorf("%s %v", registry.Join("/", s.GetPlatName(), "var", tpName, name), err)
	}
	key := fmt.Sprintf("%s/%s:%d", tpName, name, cacheConf.GetVersion())
	ok, ch, err := s.cacheMap.SetIfAbsentCb(key, func(input ...interface{}) (c interface{}, err error) {
		c, err = f(cacheConf)
		if err != nil {
			return nil, err
		}
		switch v := c.(type) {
		case CloseHandler:
			s.closeList = append(s.closeList, v)
		}
		return c, nil
	})
	if err != nil {
		err = fmt.Errorf("创建对象失败:%s,err:%v", string(cacheConf.GetRaw()), err)
		return ok, nil, err
	}
	return ok, ch, err
}

//Close 关闭缓存连接
func (s *GlobalVarObjectCache) Close() error {
	s.cacheMap.Clear()
	for _, f := range s.closeList {
		f.Close()
	}
	return nil
}
