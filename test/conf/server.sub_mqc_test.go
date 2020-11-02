package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_mqcSub_GetMQCMainConf(t *testing.T) {
	platName := "platName1"
	sysName := "sysName1"
	serverType := global.API
	clusterName := "cluster1"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	//错误的服务类型,获取配置失败
	confM := mocks.NewConfBy(platName, clusterName)
	confM.API(":8080")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
	mqcConf, err := gotS.GetMQCMainConf()
	assert.Equal(t, false, err == nil, "测试conf初始化,获取mqc对象失败")
	var nilMqc *mqc.Server
	assert.Equal(t, nilMqc, mqcConf, "测试conf初始化,判断mqc节点对象")

	//返回正确的mqc配置
	serverType = global.MQC
	confM = mocks.NewConfBy(platName, clusterName)
	confN := confM.MQC("redis://11")
	confN.Queue(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	mqcConf, err = gotS.GetMQCMainConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取mqc对象失败1")
	mqcC := mqc.New("redis://11")
	assert.Equal(t, mqcC, mqcConf, "测试conf初始化,判断mqc节点对象1")
}

func Test_mqcSub_GetMQCQueueConf(t *testing.T) {
	platName := "platName2"
	sysName := "sysName2"
	serverType := global.MQC
	clusterName := "cluster2"
	rgst, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "测试conf初始化,获取注册中心对象失败")

	confM := mocks.NewConfBy(platName, clusterName)
	confN := confM.MQC("redis://11")
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err := app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")

	//不设置queue返回结果
	queuesObj, err := gotS.GetMQCQueueConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败")
	assert.Equal(t, &queue.Queues{}, queuesObj, "测试conf初始化,判断queues节点对象")

	//设置错误的queue返回配置
	confN.Queue(queue.NewQueue("错误配置", "service1"), queue.NewQueue("错误配置1", "service2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点1")
	queuesObj, err = gotS.GetMQCQueueConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败1")
	var nilQueue *queue.Queues
	assert.Equal(t, nilQueue, queuesObj, "测试conf初始化,判断queues节点对象1")

	//设置正确的queue返回配置
	confM = mocks.NewConfBy(platName, clusterName)
	confN = confM.MQC("redis://11")
	confN.Queue(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
	gotS, err = app.NewAPPConfBy(platName, sysName, serverType, clusterName, rgst)
	assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点2")
	queuesObj, err = gotS.GetMQCQueueConf()
	assert.Equal(t, true, err == nil, "测试conf初始化,获取queues对象失败2")
	queueC := queue.NewQueues(queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2"))
	assert.Equal(t, queueC, queuesObj, "测试conf初始化,判断queues节点对象2")
}
