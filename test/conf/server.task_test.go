package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/security/md5"

	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name    string
		cron    string
		service string
		opts    []task.Option
		want    *task.Task
	}{
		{name: "空对象获取", cron: "", service: "", opts: nil, want: &task.Task{Cron: "", Service: "", Disable: false}},
		{name: "设置对象获取", cron: "cron1", service: "service1", opts: nil, want: &task.Task{Cron: "cron1", Service: "service1", Disable: false}},
		{name: "设置对象获取disable", cron: "cron1", service: "service1", opts: []task.Option{task.WithDisable()}, want: &task.Task{Cron: "cron1", Service: "service1", Disable: true}},
		{name: "设置对象获取enable", cron: "cron1", service: "service1", opts: []task.Option{task.WithEnable()}, want: &task.Task{Cron: "cron1", Service: "service1", Disable: false}},
	}
	for _, tt := range tests {
		got := task.NewTask(tt.cron, tt.service, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestTask_GetUNQ(t *testing.T) {

	tests := []struct {
		name   string
		fields *task.Task
		want   string
	}{
		{name: "空对象获取唯一标识", fields: task.NewTask("", ""), want: md5.Encrypt(fmt.Sprintf("%s(%s)", "", ""))},
		{name: "对象获取唯一标识", fields: task.NewTask("xxx", "yyy"), want: md5.Encrypt(fmt.Sprintf("%s(%s)", "yyy", "xxx"))},
	}
	for _, tt := range tests {
		got := tt.fields.GetUNQ()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestTask_IsOnce(t *testing.T) {
	tests := []struct {
		name   string
		fields *task.Task
		want   bool
	}{
		{name: "空对象执行次数判定", fields: task.NewTask("", ""), want: false},
		{name: "错误配置,对象执行次数判定", fields: task.NewTask("xxx", "yyy"), want: false},
		{name: "@once,对象执行次数判定", fields: task.NewTask("@once", "yyy"), want: true},
		{name: "@now,对象执行次数判定", fields: task.NewTask("@now", "yyy"), want: true},
	}
	for _, tt := range tests {
		got := tt.fields.IsOnce()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fields  *task.Task
		wantErr bool
	}{
		{name: "cron空报错,数据合法性判定", fields: task.NewTask("", "yyy"), wantErr: true},
		{name: "service空报错,数据合法性判定", fields: task.NewTask("xxx", ""), wantErr: true},
		{name: "cron中文报错,数据合法性判定", fields: task.NewTask("中文报错", "yyy"), wantErr: true},
		{name: "service中文报错,数据合法性判定", fields: task.NewTask("xxxx", "中文报错"), wantErr: true},
		{name: "正确数据,数据合法性判定", fields: task.NewTask("xxx", "yyy"), wantErr: false},
	}
	for _, tt := range tests {
		err := tt.fields.Validate()
		assert.Equal(t, tt.wantErr, (err != nil), tt.name)
	}
}

func TestNewEmptyTasks(t *testing.T) {
	tests := []struct {
		name string
		want *task.Tasks
	}{
		{name: "默认空对象测试", want: &task.Tasks{Tasks: make([]*task.Task, 0)}},
	}
	for _, tt := range tests {
		got := task.NewEmptyTasks()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewTasks(t *testing.T) {
	tests := []struct {
		name string
		args []*task.Task
		want *task.Tasks
	}{
		{name: "初始化入参为空nil", args: nil, want: task.NewEmptyTasks()},
		{name: "初始化入参为空", args: []*task.Task{}, want: task.NewEmptyTasks()},
		{name: "添加单个数据", args: []*task.Task{task.NewTask("xxx", "yyy")}, want: &task.Tasks{Tasks: []*task.Task{task.NewTask("xxx", "yyy")}}},
		{name: "添加多个数据", args: []*task.Task{task.NewTask("xxx", "yyy"), task.NewTask("xxx1", "yyy1")}, want: &task.Tasks{Tasks: []*task.Task{task.NewTask("xxx", "yyy"), task.NewTask("xxx1", "yyy1")}}},
	}
	for _, tt := range tests {
		got := task.NewTasks(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestTasks_Append(t *testing.T) {
	tests := []struct {
		name   string
		fields *task.Tasks
		args   []*task.Task
		want   *task.Tasks
	}{
		{name: "空对象累加对象", fields: task.NewEmptyTasks(), args: []*task.Task{task.NewTask("xxx", "yyy")}, want: task.NewTasks(task.NewTask("xxx", "yyy"))},
		{name: "空对象累加空对象", fields: task.NewEmptyTasks(), args: []*task.Task{}, want: task.NewEmptyTasks()},
		{name: "空对象累加nil对象", fields: task.NewEmptyTasks(), args: nil, want: task.NewEmptyTasks()},
		{name: "实体对象累加nil对象", fields: task.NewTasks(task.NewTask("xxx", "yyy")), args: nil, want: task.NewTasks(task.NewTask("xxx", "yyy"))},
		{name: "实体对象累加空对象", fields: task.NewTasks(task.NewTask("xxx", "yyy")), args: []*task.Task{}, want: task.NewTasks(task.NewTask("xxx", "yyy"))},
		{name: "实体对象累加单个对象", fields: task.NewTasks(task.NewTask("xxx", "yyy")), args: []*task.Task{task.NewTask("xxx1", "yyy1")}, want: task.NewTasks(task.NewTask("xxx", "yyy"), task.NewTask("xxx1", "yyy1"))},
		{name: "实体对象累加单个对象", fields: task.NewTasks(task.NewTask("xxx", "yyy")), args: []*task.Task{task.NewTask("xxx1", "yyy1"), task.NewTask("xxx2", "yyy2")}, want: task.NewTasks(task.NewTask("xxx", "yyy"), task.NewTask("xxx1", "yyy1"), task.NewTask("xxx2", "yyy2"))},
	}
	for _, tt := range tests {
		got := tt.fields.Append(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestTasksGetConf(t *testing.T) {
	type test struct {
		name    string
		cnf     conf.IServerConf
		want    *task.Tasks
		wantErr bool
	}

	conf := mocks.NewConfBy("hydra", "graytest")
	confB := conf.CRON(cron.WithTrace())
	test1 := test{name: "task节点不存在", cnf: conf.GetCronConf().GetServerConf(), want: &task.Tasks{Tasks: []*task.Task{}}, wantErr: false}
	queueObj, err := task.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name)
	if err == nil {
		assert.Equal(t, len(test1.want.Tasks), len(queueObj.Tasks), test1.name)
	}
	confB = conf.CRON(cron.WithTrace())
	confB.Task(task.NewTask("中文错误", "s2"))
	test2 := test{name: "task节点存在,数据错误", cnf: conf.GetCronConf().GetServerConf(), want: nil, wantErr: true}
	queueObj, err = task.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.want, queueObj, test2.name+",obj")

	confB = conf.CRON(cron.WithTrace())
	confB.Task(task.NewTask("@once", "s2"))
	test3 := test{name: "task节点存在,数据正确", cnf: conf.GetCronConf().GetServerConf(), want: task.NewTasks(task.NewTask("@once", "s2")), wantErr: false}
	queueObj, err = task.GetConf(test3.cnf)
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.want, queueObj, test3.name+",obj")
}
