package conf

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_cache_GetServerConf(t *testing.T) {
	//复制一个空的chache对象
	cacheobj := *server.Cache
	cacheobjPrt := &cacheobj
	//在空的缓存对象获取配置
	obj, err := cacheobjPrt.GetServerConf("api")
	assert.Equal(t, true, (err != nil), "获取不存在的配置对象，err")
	assert.Equal(t, nil, obj, "获取不存在的配置对象，res")

	conf := mocks.NewConf()
	conf.API(":8080")
	serveC := conf.GetAPIConf()
	cacheobjPrt.Save(serveC)
	obj, err = cacheobjPrt.GetServerConf("api")
	assert.Equal(t, false, (err != nil), "获取api配置对象，err")
	assert.Equal(t, serveC, obj, "获取api的配置对象，res")

	obj, err = cacheobjPrt.GetServerConf("rpc")
	assert.Equal(t, true, (err != nil), "获取rpc的配置对象，err")
	assert.Equal(t, nil, obj, "获取rpc的配置对象，res")

	serverMap := cacheobjPrt.GetServerMaps()
	key := fmt.Sprintf("%s-%v", "api", serveC.GetMainConf().GetVersion())
	serverMap.Remove(key)
	obj, err = cacheobjPrt.GetServerConf("api")
	assert.Equal(t, true, (err != nil), "获取cuur存在的配置不存在，err")
	assert.Equal(t, nil, obj, "获取cuur存在的配置不存在，res")
}

func Test_cache_GetVarConf(t *testing.T) {
	//复制一个空的chache对象
	cacheobj := *server.Cache
	cacheobjPrt := &cacheobj
	//在空的缓存对象获取配置
	obj, err := cacheobjPrt.GetVarConf()
	assert.Equal(t, true, (err != nil), "获取不存在的配置对象，err")
	assert.Equal(t, nil, obj, "获取不存在的配置对象，res")

	//却笑var设置的逻辑测试
	// conf := mocks.NewConf()
	// conf.API(":8080")

}

func Test_cache_Clear(t *testing.T) {

	conf := mocks.NewConf()
	conf.API(":8080")
	serveC := conf.GetAPIConf()
	server.Cache.Save(serveC)
	oldVerion := serveC.GetMainConf().GetVersion()
	assert.Equal(t, serveC.GetMainConf().GetVersion(), server.Cache.GetServerCuurVerion("api"), "api serverconf verion Equl")
	assert.Equal(t, oldVerion, server.Cache.GetServerCuurVerion("api"), "当前版本号比对")

	//在添加一个server配置
	serveB := conf.GetAPIConf()
	cuurtVerion := serveB.GetMainConf().GetVersion()
	server.Cache.Save(serveB)
	_, ok := server.Cache.GetServerMaps().Get(fmt.Sprintf("api-%v", oldVerion))
	assert.Equal(t, ok, true, "清除前，旧配置存在")
	assert.Equal(t, cuurtVerion, server.Cache.GetServerCuurVerion("api"), "清除前，当前版本号比对")

	time.Sleep(51 * time.Second)
	_, ok = server.Cache.GetServerMaps().Get(fmt.Sprintf("api-%v", oldVerion))
	assert.Equal(t, ok, false, "清除后，旧配置不存在")
	_, ok = server.Cache.GetServerMaps().Get(fmt.Sprintf("api-%v", cuurtVerion))
	assert.Equal(t, ok, true, "清除后，新配置存在")
	assert.Equal(t, cuurtVerion, server.Cache.GetServerCuurVerion("api"), "清除后，当前版本号比对")
}
