package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

const varNodeName = "var"

//Cache 服务器缓存信息
var Cache = &cache{
	serverMaps:           cmap.New(4),
	varMaps:              cmap.New(4),
	currentServerVersion: cmap.New(2),
	currentVarVersion:    cmap.New(2),
}

//cache通过版本号控制更新配置时引起的冲突并减少对象的拷贝
type cache struct {
	serverMaps           cmap.ConcurrentMap
	varMaps              cmap.ConcurrentMap
	currentServerVersion cmap.ConcurrentMap
	currentVarVersion    cmap.ConcurrentMap
	lock                 sync.RWMutex
}

//Save 缓存服务器配置信息
func (c *cache) Save(s IAPPConf) {
	sversion := s.GetServerConf().GetVersion()
	vversion := s.GetVarConf().GetVersion()
	typ := s.GetServerConf().GetServerType()
	c.lock.Lock()
	defer c.lock.Unlock()
	c.serverMaps.Set(getKey(typ, sversion), s)
	c.varMaps.Set(getKey(varNodeName, vversion), s.GetVarConf())
	c.currentServerVersion.Set(typ, sversion)
	c.currentVarVersion.Set(varNodeName, vversion)

}

//Get 从缓存中获取服务器配置
func (c *cache) GetAPPConf(serverType string) (IAPPConf, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	serverVerion, ok := c.currentServerVersion.Get(serverType)
	if !ok {
		return nil, fmt.Errorf("未找到%s的缓存配置信息", serverType)
	}
	if s, ok := c.serverMaps.Get(getKey(serverType, serverVerion)); ok {
		return s.(IAPPConf), nil
	}
	return nil, fmt.Errorf("获取服务器配置失败，未找到服务器[%s.%d]的缓存数据", serverType, serverVerion)
}

//Get 从缓存中获取服务器配置
func (c *cache) GetVarConf() (conf.IVarConf, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	varVerion, ok := c.currentVarVersion.Get(varNodeName)
	if !ok {
		return nil, fmt.Errorf("未找到var缓存配置信息")
	}
	if s, ok := c.varMaps.Get(getKey(varNodeName, varVerion)); ok {
		return s.(conf.IVarConf), nil
	}
	return nil, fmt.Errorf("获取var配置失败，缓存中不存在版本[%v]的数据", varVerion)
}

func (c *cache) GetServerMaps() cmap.ConcurrentMap {
	return c.serverMaps
}

func (c *cache) GetVarMaps() cmap.ConcurrentMap {
	return c.varMaps
}

func (c *cache) GetServerCuurVerion(tp string) interface{} {
	verion, ok := c.currentServerVersion.Get(tp)
	if !ok {
		return nil
	}
	return verion
}

func (c *cache) GetVarCuurVerion(varNodeName string) interface{} {
	verion, ok := c.currentVarVersion.Get(varNodeName)
	if !ok {
		return nil
	}
	return verion.(int32)
}
func (c *cache) clear() {
	tm := time.NewTicker(time.Second * 50)
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-tm.C:

			c.serverMaps.RemoveIterCb(func(key string, v interface{}) bool {
				conf := v.(IAPPConf)
				tp := conf.GetServerConf().GetServerType()
				ver, _ := c.currentServerVersion.Get(tp)
				currentKey := getKey(tp, ver)

				if key != currentKey {
					conf.Close()
					global.Def.Log().Debug("清理缓存配置[%s]", conf.GetServerConf().GetServerPath())
					return true
				}
				return false
			})
			c.varMaps.RemoveIterCb(func(key string, v interface{}) bool {
				currentVer, _ := c.currentVarVersion.Get(varNodeName)
				currentKey := getKey(varNodeName, currentVer)
				if key != currentKey {
					global.Def.Log().Debug("清理缓存配置[%s]", key)
					return true
				}
				return false
			})
		}
	}
}
func getKey(tp string, version interface{}) string {
	return fmt.Sprintf("%s-%v", tp, version)
}
func init() {
	go Cache.clear()
}
