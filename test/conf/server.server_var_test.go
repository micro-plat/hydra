package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/registry"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewVarConf(t *testing.T) {
	systemName := "sys1"
	platName := "hydra1"
	clusterName := "cluter1"
	confM := mocks.NewConfBy(platName, clusterName)
	confM.Vars()
	confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)

	//错误的var路径初始化空对象
	varConf, err := server.NewVarConf("varConfPath", confM.Registry)
	assert.Equal(t, true, err == nil, "初始化varconf失败")
	dbConf, err := varConf.GetConf("db", "db")
	assert.Equal(t, conf.ErrNoSetting, err, "获取db节点配置异常")
	assert.Equal(t, conf.EmptyJSONConf, dbConf, "获取db节点配置数据不正确")
	assert.Equal(t, int32(0), varConf.GetVersion(), "获取的版本号是错误的")

	vsion, err := varConf.GetConfVersion("db", "db")
	assert.Equal(t, conf.ErrNoSetting, err, "获取db子节点版本号异常")
	assert.Equal(t, int32(0), vsion, "获取db子节点版本号不正确")

	dbObj := db.DB{}
	vsion, err = varConf.GetObject("db", "db", &dbObj)
	assert.Equal(t, conf.ErrNoSetting, err, "获取db子节点对象版本号异常")
	assert.Equal(t, int32(0), vsion, "获取db子节点对象版本号不正确")
	assert.Equal(t, db.DB{}, dbObj, "获取db子节点对象失败")

	varConf1 := varConf.GetClone()
	assert.Equal(t, varConf1, varConf, "克隆的var节点不正确")

	assert.Equal(t, false, varConf.Has("db", "db"), "var节点是否存在判断错误")

	//正确的var路径初始化空配置对象
	varConfPath := server.NewVarPath(platName).GetVarPath()
	varConf, err = server.NewVarConf(varConfPath, confM.Registry)
	assert.Equal(t, true, err == nil, "初始化varconf失败1")
	dbConf, err = varConf.GetConf("db", "db")
	assert.Equal(t, conf.ErrNoSetting, err, "获取db节点配置异常1")
	assert.Equal(t, conf.EmptyJSONConf, dbConf, "获取db节点配置数据不正确1")

	vsion, err = varConf.GetConfVersion("db", "db")
	assert.Equal(t, conf.ErrNoSetting, err, "获取db子节点版本号异常1")
	assert.Equal(t, int32(0), vsion, "获取db子节点版本号不正确1")

	dbObj = db.DB{}
	vsion, err = varConf.GetObject("db", "db", &dbObj)
	assert.Equal(t, conf.ErrNoSetting, err, "获取db子节点对象版本号异常1")
	assert.Equal(t, int32(0), vsion, "获取db子节点对象版本号不正确1")
	assert.Equal(t, db.DB{}, dbObj, "获取db子节点对象失败1")

	varConf1 = varConf.GetClone()
	assert.Equal(t, varConf1, varConf, "克隆的var节点不正确1")
	assert.Equal(t, false, varConf.Has("db", "db"), "var节点是否存在判断错误1")

	//设置全新的db节点var配置对象
	systemName = "sys2"
	platName = "hydra2"
	clusterName = "cluter2"
	confN := mocks.NewConfBy(platName, clusterName)
	confN.Vars().DB("newdb", oracle.NewBy("taosy", "123456", "tnsName"))
	confN.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
	varConfPath = server.NewVarPath(platName).GetVarPath()
	varConf, err = server.NewVarConf(varConfPath, confN.Registry)
	assert.Equal(t, true, err == nil, "初始化varconf失败2")
	dbConf, err = varConf.GetConf("db", "newdb")
	assert.Equal(t, true, err == nil, "获取db节点配置异常2")
	assert.Equal(t, "oracle", dbConf.GetString("provider"), "获取db节点Provider配置数据不正确")
	assert.Equal(t, "taosy/123456@tnsName", dbConf.GetString("connString"), "获取db节点connString配置数据不正确")
	assert.Equal(t, 10, dbConf.GetInt("maxOpen"), "获取db节点maxOpen配置数据不正确")
	assert.Equal(t, 3, dbConf.GetInt("maxIdle"), "获取db节点maxIdle配置数据不正确")
	assert.Equal(t, 600, dbConf.GetInt("lifeTime"), "获取db节点lifeTime配置数据不正确")

	_, vsion, err = confN.Registry.GetValue(varConfPath)
	assert.Equal(t, true, err == nil, "注册中心获取节点数据异常2")
	assert.Equal(t, vsion, varConf.GetVersion(), "获取的版本号是错误的2")

	_, vsion1, err := confN.Registry.GetValue(registry.Join(varConfPath, "db", "newdb"))
	assert.Equal(t, true, err == nil, "注册中心获取节点数据异常3")
	vsion, err = varConf.GetConfVersion("db", "newdb")
	assert.Equal(t, true, err == nil, "获取db子节点版本号异常2")
	assert.Equal(t, vsion1, vsion, "获取db子节点版本号不正确3")

	dbObj = db.DB{}
	vsion, err = varConf.GetObject("db", "newdb", &dbObj)
	assert.Equal(t, true, err == nil, "获取db子节点对象版本号异常2")
	assert.Equal(t, vsion1, vsion, "获取db子节点对象版本号不正确2")
	assert.Equal(t, "oracle", dbObj.Provider, "获取db子节点对象失败,Provider")
	assert.Equal(t, 10, dbObj.MaxOpen, "获取db子节点对象失败,MaxOpen")
	assert.Equal(t, 3, dbObj.MaxIdle, "获取db子节点对象失败,MaxIdle")
	assert.Equal(t, 600, dbObj.LifeTime, "获取db子节点对象失败,LifeTime")
	assert.Equal(t, "taosy/123456@tnsName", dbObj.ConnString, "获取db子节点对象失败,ConnString")

	varConf1 = varConf.GetClone()
	assert.Equal(t, varConf1, varConf, "克隆的var节点不正确2")
	assert.Equal(t, true, varConf.Has("db", "newdb"), "var节点是否存在判断错误2")

}
