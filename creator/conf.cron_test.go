package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_newCron(t *testing.T) {
	tests := []struct {
		name string
		opts []cron.Option
		want *cronBuilder
	}{
		{name: "生成空对象", opts: []cron.Option{}, want: &cronBuilder{CustomerBuilder: map[string]interface{}{"main": cron.New()}}},
		{name: "生成实体对象对象", opts: []cron.Option{cron.WithDisable(), cron.WithMasterSlave()}, want: &cronBuilder{CustomerBuilder: map[string]interface{}{"main": cron.New(cron.WithDisable(), cron.WithMasterSlave())}}},
	}
	for _, tt := range tests {
		got := newCron(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_cronBuilder_Load(t *testing.T) {

	tasks1, _ := task.NewEmptyTasks().Append(task.NewTask("cron", "server"))
	tasks2, _ := task.NewEmptyTasks().Append(task.NewTask("cron", "server"), task.NewTask("cron1", "server1"))
	tests := []struct {
		name    string
		addlist map[string]string
		obj     *cronBuilder
		want    *cronBuilder
	}{
		{name: "全空数据", addlist: map[string]string{}, obj: &cronBuilder{CustomerBuilder: map[string]interface{}{}}, want: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": task.NewEmptyTasks()}}},
		{name: "实体对象,add任务", addlist: map[string]string{"cron1": "server1"}, obj: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": tasks1}},
			want: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": tasks2}}},
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
		{name: "全空数据", tks: []*task.Task{}, obj: &cronBuilder{CustomerBuilder: map[string]interface{}{}}, want: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": task.NewTasks()}}},
		{name: "空数据,加对象", tks: []*task.Task{task.NewTask("cron", "service")}, obj: &cronBuilder{CustomerBuilder: map[string]interface{}{}},
			want: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"))}}},
		{name: "有数据,加对象", tks: []*task.Task{task.NewTask("cron1", "service1")}, obj: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"))}},
			want: &cronBuilder{CustomerBuilder: map[string]interface{}{"task": task.NewTasks(task.NewTask("cron", "service"), task.NewTask("cron1", "service1"))}}},
	}
	for _, tt := range tests {
		tt.obj.Task(tt.tks...)
		assert.Equal(t, tt.want, tt.obj, tt.name)
	}
}
