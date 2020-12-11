package creator

import (
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rpc"

	"github.com/micro-plat/hydra/conf/vars/http"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

type vars map[string]map[string]interface{}

//DB 添加db配置
func (v vars) DB() *Vardb {
	return NewDB(v)
}

func (v vars) Cache() *Varcache {
	return NewCache(v)
}

func (v vars) Queue() *Varqueue {
	return NewQueue(v)
}

func (v vars) RLog(service string, opts ...rlog.Option) vars {
	v.Custom(rlog.TypeNodeName, rlog.LogName, rlog.New(service, opts...))
	return v
}

func (v vars) HTTP(nodeName string, opts ...http.Option) vars {
	v.Custom(http.HttpTypeNode, nodeName, http.New(opts...))
	return v
}

func (v vars) RPC(nodeName string, opts ...rpc.Option) vars {
	v.Custom(rpc.RPCTypeNode, nodeName, rpc.New(opts...))
	return v
}

//Redis 添加Redis配置
func (v vars) Redis(nodeName string, address string, opts ...redis.Option) vars {
	v.Custom(redis.TypeNodeName, nodeName, redis.New(address, opts...))
	return v
}

//Custom 自定义配置
func (v vars) Custom(typ string, nodeName string, i interface{}) vars {
	if _, ok := v[typ]; !ok {
		v[typ] = make(map[string]interface{})
	}
	v[typ][nodeName] = i
	return v
}
