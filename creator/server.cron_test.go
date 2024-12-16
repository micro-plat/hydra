package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/assert"
)

func Test_newCron(t *testing.T) {
	tests := []struct {
		name string
		opts []cron.Option
		want *cronBuilder
	}{
		{name: "1. 初始化空对象", opts: []cron.Option{}, want: &cronBuilder{BaseBuilder: map[string]interface{}{"main": cron.New()}}},
		{name: "2. 初始化实体对象", opts: []cron.Option{cron.WithDisable(), cron.WithMasterSlave()}, want: &cronBuilder{BaseBuilder: map[string]interface{}{"main": cron.New(cron.WithDisable(), cron.WithMasterSlave())}}},
	}
	for _, tt := range tests {
		got := newCron(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_cronBuilder_Load(t *testing.T) {

	tasks1, _ := task.NewEmptyTasks().Append(task.NewTask("cron", "server"))
	tasks2, _ := task.NewEmptyTasks().Append(task.NewTask("cron", "server"), task.NewTask("cron1", "server1"))
	tasks3 := &task.Tasks{Tasks: []*task.Task{task.NewTask("cron", "server"), task.NewTask("cron1", "server1")}}
	tests := []struct {
		name    string
		addlist map[string]string
		obj     *cronBuilder
		want    *cronBuilder
	}{
		{name: "1. 加载空任务对象", addlist: map[string]string{}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{}}, want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": task.NewEmptyTasks()}}},
		{name: "2. 加载实体对象", addlist: map[string]string{"cron1": "server1"}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{"task": tasks1}},
			want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": tasks1}}},
		{name: "3. 加载重复的任务对象", addlist: map[string]string{"cron1": "server1"}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{"task": tasks2}},
			want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": tasks3}}},
	}
	for _, tt := range tests {
		for k, v := range tt.addlist {
			services.CRON.Add(k, v)
		}
		tt.obj.Load()
		assert.Equal(t, tt.want, tt.obj, tt.name)
	}
}

func Test_cronBuilder_Task(t *testing.T) {
	tests := []struct {
		name string
		tks  []*task.Task
		obj  *cronBuilder
		want *cronBuilder
	}{
		{name: "1. 全空数据", tks: []*task.Task{}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{}}, want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": task.NewTasks()}}},
		{name: "2. 空对象添加任务数据", tks: []*task.Task{task.NewTask("cron", "service")}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{}},
			want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"))}}},
		{name: "3. 有数据对象累加任务数据", tks: []*task.Task{task.NewTask("cron1", "service1")}, obj: &cronBuilder{BaseBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"))}},
			want: &cronBuilder{BaseBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"), task.NewTask("cron1", "service1"))}}},
	}
	for _, tt := range tests {
		tt.obj.Task(tt.tks...)
		assert.Equal(t, tt.want, tt.obj, tt.name)
	}
}
