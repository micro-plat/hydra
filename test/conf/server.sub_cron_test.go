package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_cronSub_GetCRONTaskConf(t *testing.T) {
	platName := "platName1"
	sysName := "sysName1"
	serverType := global.CRON
	clusterName := "cluster1"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	confN := confM.CRON(cron.WithTrace(), cron.WithDisable(), cron.WithSharding(1))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	//空task获取对象
	taskConf, err := gotS.GetCRONTaskConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败")
	assert.Equal(t, &task.Tasks{Tasks: []*task.Task{}}, taskConf, "测试conf初始化,判断task节点对象")

	//设置错误数据的task获取对象
	confM = mocks.NewConfBy(platName, clusterName)
	confN = confM.CRON()
	confN.Task(task.NewTask("错误数据", "service1"), task.NewTask("错误数据1", "service2", task.WithDisable()))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	taskConf, err = gotS.GetCRONTaskConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取cron对象失败1")
	var nilTask *task.Tasks
	assert.Equal(t, nilTask, taskConf, "测试conf初始化,判断task节点对象1")

	//设置正确的task对象
	confM = mocks.NewConfBy(platName, clusterName)
	confN = confM.CRON()
	confN.Task(task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable()))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	taskConf, err = gotS.GetCRONTaskConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取cron对象失败2")
	taskC := task.NewTasks(task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable()))
	assert.Equal(t, taskC, taskConf, "测试conf初始化,判断task节点对象2")
}
