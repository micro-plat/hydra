package builder

import (
	"github.com/micro-plat/hydra/registry/conf/server/cron"
)

type cronBuilder map[string]interface{}

//newCron 构建cron生成器
func newCron(opts ...cron.Option) cronBuilder {
	b := make(map[string]interface{})
	b["main"] = cron.New(opts...)
	return b
}
