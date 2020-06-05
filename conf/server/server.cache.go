package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//Cache 服务器缓存信息
var Cache = &cache{
	serverMaps:     cmap.New(4),
	varMaps:        cmap.New(4),
	versionHistory: make([]int32, 0, 10),
}

//cache通过版本号控制更新配置时引起的冲突并减少对象的拷贝
type cache struct {
	serverMaps           cmap.ConcurrentMap
	varMaps              cmap.ConcurrentMap
	versionHistory       []int32
	currentServerVersion int32
	currentVarVersion    int32
	lock                 sync.RWMutex
}

//Save 缓存服务器配置信息
func (c *cache) Save(s IServerConf) {
	sversion := s.GetMainConf().GetVersion()
	vversion := s.GetVarConf().GetVersion()
	typ := s.GetMainConf().GetServerType()
	key := fmt.Sprintf("%s_%d", typ, sversion)
	c.lock.Lock()
	defer c.lock.Unlock()
	c.serverMaps.Set(key, s)
	c.varMaps.Set(fmt.Sprint(vversion), s.GetVarConf())
	c.currentServerVersion = sversion
	c.currentVarVersion = vversion
	c.versionHistory = append(c.versionHistory, c.currentServerVersion)
	c.versionHistory = append(c.versionHistory, c.currentVarVersion)

}

//Get 从缓存中获取服务器配置
func (c *cache) GetServerConf(serverType string) (IServerConf, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	key := fmt.Sprintf("%s_%d", serverType, c.currentServerVersion)
	if s, ok := c.serverMaps.Get(key); ok {
		return s.(IServerConf), nil
	}
	return nil, fmt.Errorf("获取服务器配置失败，未找到服务器[%s.%d]的缓存数据", serverType, c.currentServerVersion)
}

//Get 从缓存中获取服务器配置
func (c *cache) GetVarConf() (conf.IVarConf, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	key := fmt.Sprintf("%d", c.currentVarVersion)
	if s, ok := c.varMaps.Get(key); ok {
		return s.(conf.IVarConf), nil
	}
	return nil, fmt.Errorf("获取var配置失败，缓存中不存在版本[%d]的数据", c.currentVarVersion)
}

func (c *cache) clear() {
	tm := time.NewTicker(time.Second * 3600)
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-tm.C:
			c.serverMaps.RemoveIterCb(func(key string, v interface{}) bool {
				if key != fmt.Sprint(c.currentServerVersion) {
					if s, ok := v.(IServerConf); ok {
						s.Close()
					}
					global.Def.Log().Debug("清理ServerConf缓存配置", key)
					return true
				}
				return false
			})
			c.varMaps.RemoveIterCb(func(key string, v interface{}) bool {
				global.Def.Log().Debug("清理VarConf缓存配置", key)
				return key != fmt.Sprint(c.currentVarVersion)
			})
		}
	}
}
func init() {
	go Cache.clear()
}
