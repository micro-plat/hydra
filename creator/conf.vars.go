package creator

import (
	"github.com/micro-plat/hydra/conf/plat/cache"
	"github.com/micro-plat/hydra/conf/plat/db"
	"github.com/micro-plat/hydra/conf/plat/queue"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) DB(name string, db *db.DB) vars {
	if _, ok := v["db"]; !ok {
		v["db"] = make(map[string]interface{})
	}
	v["db"][name] = db
	return v
}

func (v vars) Queue(name string, q *queue.Queue) vars {
	if _, ok := v["queue"]; !ok {
		v["queue"] = make(map[string]interface{})
	}
	v["queue"][name] = q
	return v
}

func (v vars) Cache(name string, q *cache.Cache) vars {
	if _, ok := v["cache"]; !ok {
		v["cache"] = make(map[string]interface{})
	}
	v["cache"][name] = q
	return v
}
