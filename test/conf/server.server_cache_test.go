package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_cache_GetServerConf(t *testing.T) {
	//复制一个空的chache对象
	cacheobj := *app.Cache
	cacheobjPrt := &cacheobj
	//在空的缓存对象获取配置
	obj, err := cacheobjPrt.GetAPPConf("api")
	assert.Equal(t, true, (err != nil), "获取不存在的配置对象，err")
	assert.Equal(t, nil, obj, "获取不存在的配置对象，res")

	conf := mocks.NewConfBy("hydraconf_servervavhe_test2", "servercache")
	conf.API(":8080")
	serveC := conf.GetAPIConf()
	cacheobjPrt.Save(serveC)
	obj, err = cacheobjPrt.GetAPPConf("api")
	assert.Equal(t, false, (err != nil), "获取api配置对象，err")
	assert.Equal(t, serveC, obj, "获取api的配置对象，res")

	obj, err = cacheobjPrt.GetAPPConf("rpc")
	assert.Equal(t, true, (err != nil), "获取rpc的配置对象，err")
	assert.Equal(t, nil, obj, "获取rpc的配置对象，res")

	serverMap := cacheobjPrt.GetServerHistory()
	key := fmt.Sprintf("%s-%v", "api", serveC.GetServerConf().GetVersion())
	serverMap.Remove(key)
	obj, err = cacheobjPrt.GetAPPConf("api")
	assert.Equal(t, true, (err != nil), "获取cuur存在的配置不存在，err")
	assert.Equal(t, nil, obj, "获取cuur存在的配置不存在，res")
}

func Test_cache_GetVarConf(t *testing.T) {
	//复制一个空的chache对象
	cacheobj := *app.Cache
	cacheobjPrt := &cacheobj
	app.Cache.GetVarHistory().Clear()
	//在空的缓存对象获取配置
	obj, err := cacheobjPrt.GetVarConf()
	assert.Equal(t, true, (err != nil), "获取不存在的配置对象，err")
	assert.Equal(t, nil, obj, "获取不存在的配置对象，res")
}

func Test_cache_Clear(t *testing.T) {

	conf := mocks.NewConfBy("hydraconf_servervavhe_clear", "servervavhe_clear")
	conf.API(":8080")
	serveC := conf.GetAPIConf()
	app.Cache.Save(serveC)
	oldVerion := serveC.GetServerConf().GetVersion()
	assert.Equal(t, serveC.GetServerConf().GetVersion(), app.Cache.GetCurrentServerVerion("api"), "api serverconf verion Equl")
	assert.Equal(t, oldVerion, app.Cache.GetCurrentServerVerion("api"), "当前版本号比对")

	//在添加一个server配置
	serveB := conf.GetAPIConf()
	cuurtVerion := serveB.GetServerConf().GetVersion()
	app.Cache.Save(serveB)
	_, ok := app.Cache.GetServerHistory().Get(fmt.Sprintf("api-%v", oldVerion))
	assert.Equal(t, ok, true, "清除前，旧配置存在")
	assert.Equal(t, cuurtVerion, app.Cache.GetCurrentServerVerion("api"), "清除前，当前版本号比对")

	//@todo clear 的时间是5min，已超过自动测试用力的执行时间
	// time.Sleep(51 * time.Second)
	// _, ok = app.Cache.GetServerHistory().Get(fmt.Sprintf("api-%v", oldVerion))
	// assert.Equal(t, ok, false, "清除后，旧配置不存在")
	_, ok = app.Cache.GetServerHistory().Get(fmt.Sprintf("api-%v", cuurtVerion))
	assert.Equal(t, ok, true, "清除后，新配置存在")
	assert.Equal(t, cuurtVerion, app.Cache.GetCurrentServerVerion("api"), "清除后，当前版本号比对")
}
