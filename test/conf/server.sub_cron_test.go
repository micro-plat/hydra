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

	var nilTask *task.Tasks
	tests := []struct {
		name     string
		opts     []*task.Task
		wantErr  bool
		wantConf *task.Tasks
	}{
		{name: "空task获取对象", opts: []*task.Task{}, wantErr: true, wantConf: &task.Tasks{Tasks: []*task.Task{}}},
		{name: "设置错误数据的task获取对象", opts: []*task.Task{task.NewTask("错误数据", "service1"), task.NewTask("错误数据1", "service2", task.WithDisable())}, wantErr: false, wantConf: nilTask},
		{name: "设置正确的task对象", opts: []*task.Task{task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable())}, wantErr: true,
			wantConf: task.NewTasks(task.NewTask("cron1", "service1"), task.NewTask("cron2", "service2", task.WithDisable()))},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.CRON(cron.WithTrace(), cron.WithDisable(), cron.WithSharding(1))
		if len(tt.opts) > 0 {
			confN.Task(tt.opts...)
		}
		confM.Conf().Pub(platName, sysName, clusterName, "lm://.", true)
		gotS, err := server.NewServerConfBy(platName, sysName, serverType, clusterName, rgst)
		assert.Equal(t, true, err == nil, "测试conf初始化,设置主节点")
		taskConf, err := gotS.GetCRONTaskConf()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.wantConf, taskConf, tt.name+",conf")
	}
}
