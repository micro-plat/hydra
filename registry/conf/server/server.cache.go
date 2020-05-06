package server

import "github.com/micro-plat/lib4go/concurrent/cmap"

var serverMaps = cmap.New(6)

//Cache 缓存服务器配置信息
func Cache(s IServerConf) {
	serverMaps.Set(s.GetMainConf().GetServerType(), s)
}

//Get 从缓存中获取服务器配置
func Get(serverType string) (IServerConf, bool) {
	if s, ok := serverMaps.Get(serverType); ok {
		return s.(IServerConf), true
	}
	return nil, false
}
