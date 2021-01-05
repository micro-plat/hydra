package creator

import (
	"github.com/micro-plat/hydra/conf/server/rpc"
)

type rpcBuilder struct {
	*httpBuilder
}

//newHTTP 构建http生成器
func newRPC(address string, opts ...rpc.Option) *rpcBuilder {
	b := &rpcBuilder{
		httpBuilder: &httpBuilder{
			BaseBuilder: make(map[string]interface{}),
		},
	}
	b.BaseBuilder[ServerMainNodeName] = rpc.New(address, opts...)
	return b
}

//Load 加载路由
func (b *rpcBuilder) Load() {
	return
}
