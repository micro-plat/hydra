package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_varSub_GetRLogConf(t *testing.T) {
	platName := "platName1"
	sysName := "sysName1"
	serverType := global.API
	clusterName := "cluster1"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	//空Rlog获取对象
	rlogConf, err := gotS.GetRLogConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取layout对象失败")
	assert.Equal(t, &rlog.Layout{Layout: rlog.DefaultLayout, Disable: true}, rlogConf, "测试conf初始化,判断layout节点对象")

	//设置错误数据的task获取对象(设置不了错误的参数)

	//设置正确的task对象
	confM.Vars().RLog("test1")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	rlogConf, err = gotS.GetRLogConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败2")
	rlogC := rlog.New("test1")
	assert.Equal(t, rlogC, rlogConf, "测试conf初始化,判断task节点对象2")
}
