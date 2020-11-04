package creator

import (
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/creator/confvars"

	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) DB() *confvars.Vardb {
	return confvars.NewDB(v)
}

func (v vars) Cache() *confvars.Varcache {
	return confvars.NewCache(v)
}

func (v vars) Queue() *confvars.Varqueue {
	return confvars.NewQueue(v)
}

func (v vars) RLog(service string, opts ...rlog.Option) vars {
	if _, ok := v[rlog.TypeNodeName]; !ok {
		v[rlog.TypeNodeName] = make(map[string]interface{})
	}
	v[rlog.TypeNodeName][rlog.LogName] = rlog.New(service, opts...)
	return v
}

//Redis 添加Redis配置
func (v vars) Redis(name string, opts *redis.Redis) vars {
	if _, ok := v[redis.TypeNodeName]; !ok {
		v[redis.TypeNodeName] = make(map[string]interface{})
	}
	v[redis.TypeNodeName][name] = opts
	return v
}
