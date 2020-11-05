package creator

import (
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/services"
)

type cronBuilder struct {
	CustomerBuilder
}

//newCron 构建cron生成器
func newCron(opts ...cron.Option) *cronBuilder {
	b := &cronBuilder{
		CustomerBuilder: make(map[string]interface{}),
	}
	b.CustomerBuilder["main"] = cron.New(opts...)
	return b
}

//Load 加载路由
func (b *cronBuilder) Load() {
	tasks := services.CRON.GetTasks()
	if q, ok := b.CustomerBuilder["task"].(*task.Tasks); ok {
		q.Append(tasks.Tasks...)
		return
	}
	b.CustomerBuilder["task"] = tasks
	return
}

//Queue 添加队列配置
func (b *cronBuilder) Task(tks ...*task.Task) *cronBuilder {
	otask, ok := b.CustomerBuilder["task"].(*task.Tasks)
	if !ok {
		otask = task.NewTasks()
		b.CustomerBuilder["task"] = otask
	}
	otask.Append(tks...)
	return b
}
