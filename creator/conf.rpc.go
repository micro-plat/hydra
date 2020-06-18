package creator

import (
	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/services"
)

type rpcBuilder struct {
	*httpBuilder
}

//newHTTP 构建http生成器
func newRPC(address string, opts ...rpc.Option) *rpcBuilder {
	b := &rpcBuilder{
		httpBuilder: &httpBuilder{
			customerBuilder: make(map[string]interface{}),
		},
	}
	b.customerBuilder["main"] = rpc.New(address, opts...)
	return b
}

//Load 加载路由
func (b *rpcBuilder) Load() {
	routers, err := services.RPC.GetRouters()
	if err != nil {
		panic(err)
	}
	b.customerBuilder["router"] = routers
	return
}
