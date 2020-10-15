package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestRpcGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	wantC := rpc.New(":8088", rpc.WithTrace())
	conf.RPC(":8088", rpc.WithTrace())
	got, err := rpc.GetConf(conf.GetRPCConf().GetMainConf())
	if err != nil {
		t.Errorf("rpcGetConf 获取配置对对象失败,err: %v", err)
	}
	assert.Equal(t, got, wantC, "检查对象是否满足预期")
}

func TestRpcGetConf1(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	conf.RPC(":8088")
	_, err := rpc.GetConf(conf.GetRPCConf().GetMainConf())
	if err == nil {
		t.Errorf("rpcGetConf 获取怕配置不能成功")
	}
}
