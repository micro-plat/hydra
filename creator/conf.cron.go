package creator

import (
	"github.com/micro-plat/hydra/conf/server/cron"
)

type cronBuilder customerBuilder

//newCron 构建cron生成器
func newCron(opts ...cron.Option) cronBuilder {
	b := make(map[string]interface{})
	b["main"] = cron.New(opts...)
	return b
}
