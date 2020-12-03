package creator

import (
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/hydra/creator/internal"

	"github.com/micro-plat/hydra/conf/vars/http"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) DB() *internal.Vardb {
	return internal.NewDB(v)
}

func (v vars) Cache() *internal.Varcache {
	return internal.NewCache(v)
}

func (v vars) Queue() *internal.Varqueue {
	return internal.NewQueue(v)
}

func (v vars) RLog(service string, opts ...rlog.Option) vars {
	if _, ok := v[rlog.TypeNodeName]; !ok {
		v[rlog.TypeNodeName] = make(map[string]interface{})
	}
	v[rlog.TypeNodeName][rlog.LogName] = rlog.New(service, opts...)
	return v
}

func (v vars) HTTP(name string, opts ...http.Option) vars {
	if _, ok := v[http.HttpTypeNode]; !ok {
		v[http.HttpTypeNode] = make(map[string]interface{})
	}
	v[http.HttpTypeNode][name] = http.New(opts...)
	return v
}

func (v vars) RPC(name string, opts ...rpc.Option) vars {
	if _, ok := v[rpc.RPCTypeNode]; !ok {
		v[rpc.RPCTypeNode] = make(map[string]interface{})
	}
	v[rpc.RPCTypeNode][name] = rpc.New(opts...)
	return v
}

//Redis 添加Redis配置
func (v vars) Redis(name string, r *redis.Redis) vars {
	if _, ok := v[redis.TypeNodeName]; !ok {
		v[redis.TypeNodeName] = make(map[string]interface{})
	}
	v[redis.TypeNodeName][name] = r
	return v
}

//Custom 自定义配置
func (v vars) Custom(typ string, name string, i interface{}) vars {
	if _, ok := v[typ]; !ok {
		v[typ] = make(map[string]interface{})
	}
	v[typ][name] = i
	return v
}
