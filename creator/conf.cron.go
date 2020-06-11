package creator

import (
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/services"
)

type cronBuilder struct {
	customerBuilder
}

//newCron 构建cron生成器
func newCron(opts ...cron.Option) *cronBuilder {
	b := &cronBuilder{
		customerBuilder: make(map[string]interface{}),
	}
	b.customerBuilder["main"] = cron.New(opts...)
	return b
}

//Load 加载路由
func (b *cronBuilder) Load() {
	tasks := services.CRON.GetTasks()
	b.customerBuilder["task"] = tasks
	return
}
