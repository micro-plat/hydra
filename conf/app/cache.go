package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

const varNodeName = "var"

//Cache 服务器缓存信息
var Cache = &cache{
	serverConfHistory:    cmap.New(4),
	varConfHistory:       cmap.New(4),
	currentServerVersion: cmap.New(2),
	currentVarVersion:    cmap.New(2),
}

//cache通过版本号控制更新配置时引起的冲突并减少对象的拷贝
type cache struct {
	serverConfHistory    cmap.ConcurrentMap
	varConfHistory       cmap.ConcurrentMap
	currentServerVersion cmap.ConcurrentMap
	currentVarVersion    cmap.ConcurrentMap
	lock                 sync.Mutex
}

//Save 将应用配置信息存入缓存
func (c *cache) Save(s IAPPConf) {

	//获取版本号
	sVersion := s.GetServerConf().GetVersion()
	vVersion := s.GetVarConf().GetVersion()
	typ := s.GetServerConf().GetServerType()

	//控制数据不一致
	c.lock.Lock()
	defer c.lock.Unlock()

	//保存最新版本号
	c.currentServerVersion.Set(typ, sVersion)
	c.currentVarVersion.Set(varNodeName, vVersion)

	//保存配置信息
	c.serverConfHistory.Set(getKey(typ, sVersion), s)
	c.varConfHistory.Set(getKey(varNodeName, vVersion), s.GetVarConf())
}

//GetAPPConf 根据服务器类型获取配置
func (c *cache) GetAPPConf(serverType string) (IAPPConf, error) {

	//获取配置版本号
	serverVerion, ok := c.currentServerVersion.Get(serverType)
	if !ok {
		return nil, fmt.Errorf("未找到%s的缓存配置信息", serverType)
	}

	//获取配置信息
	if s, ok := c.serverConfHistory.Get(getKey(serverType, serverVerion)); ok {
		return s.(IAPPConf), nil
	}
	return nil, fmt.Errorf("获取服务器配置失败，未找到服务器[%s.%d]的缓存数据", serverType, serverVerion)
}

//GetVarConf 从缓存中获取var配置
func (c *cache) GetVarConf() (conf.IVarConf, error) {

	//获取配置版本号
	varVerion, ok := c.currentVarVersion.Get(varNodeName)
	if !ok {
		return nil, fmt.Errorf("未找到var缓存配置信息,或配置信息未准备好")
	}

	//获取配置信息
	if s, ok := c.varConfHistory.Get(getKey(varNodeName, varVerion)); ok {
		return s.(conf.IVarConf), nil
	}
	return nil, fmt.Errorf("获取var配置失败，缓存中不存在版本[%v]的数据", varVerion)
}

//GetServerHistory 获取服务器配置历史
func (c *cache) GetServerHistory() cmap.ConcurrentMap {
	return c.serverConfHistory
}

//GetVarHistory 获取var服务器配置历史
func (c *cache) GetVarHistory() cmap.ConcurrentMap {
	return c.varConfHistory
}

//GetCurrentServerVerion 服务器可用版本号
func (c *cache) GetCurrentServerVerion(tp string) int32 {
	verion, ok := c.currentServerVersion.Get(tp)
	if !ok {
		return 0
	}
	return verion.(int32)
}

//GetCurrentVarVerion 获取var版本号
func (c *cache) GetCurrentVarVerion(varNodeName string) int32 {
	verion, ok := c.currentVarVersion.Get(varNodeName)
	if !ok {
		return 0
	}
	return verion.(int32)
}

//clear 定时清理历史配置
func (c *cache) clear() {
	tm := time.NewTicker(time.Minute * 5)
LOOP:
	for {
		select {
		case <-global.Def.ClosingNotify():
			break LOOP
		case <-tm.C:

			//清理服务器配置
			c.serverConfHistory.RemoveIterCb(func(key string, v interface{}) bool {
				conf := v.(IAPPConf)
				tp := conf.GetServerConf().GetServerType()
				ver, _ := c.currentServerVersion.Get(tp)
				currentKey := getKey(tp, ver)
				if key != currentKey {
					conf.Close()
					global.Def.Log().Debugf("清理缓存配置[%s]", conf.GetServerConf().GetServerPath())
					return true
				}
				return false
			})

			//清理var配置
			c.varConfHistory.RemoveIterCb(func(key string, v interface{}) bool {
				currentVer, _ := c.currentVarVersion.Get(varNodeName)
				currentKey := getKey(varNodeName, currentVer)
				if key != currentKey {
					global.Def.Log().Debugf("清理缓存配置[%s]", key)
					return true
				}
				return false
			})
		}
	}
}

//getKey 获取缓存key
func getKey(tp string, version interface{}) string {
	return fmt.Sprintf("%s-%v", tp, version)
}

//PullAndSave 拉取注册中心的服务器配置，并缓存
func PullAndSave() error {
	//接取配置信息
	for _, tp := range global.Def.ServerTypes {
		pub := server.NewServerPub(global.Def.GetPlatName(), global.Def.SysName, tp, global.Def.ClusterName)
		conf, err := NewAPPConf(pub.GetServerPath(), registry.GetCurrent())
		if err != nil {
			return fmt.Errorf("获取%s配置发生错误:%v", pub.GetServerPath(), err)
		}
		//保存配置缓存
		Cache.Save(conf)
	}
	return nil
}

//定时清理缓存
func init() {
	go Cache.clear()
}
