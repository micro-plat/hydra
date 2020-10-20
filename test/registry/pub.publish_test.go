package registry

import (
	"testing"

	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestPublisher_pubServerNode(t *testing.T) {
	type args struct {
		serverName  string
		serviceAddr string
		clusterID   string
		service     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "发布服务", args: args{serverName: "192.168.5.115:9999", serviceAddr: "lm://.", clusterID: "", service: []string{}}},
	}
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置
	p := pub.New(c)
	for _, tt := range tests {
		err := p.Publish(tt.args.serverName, tt.args.serviceAddr, c.GetServerID(), s.GetRouterConf().GetPath()...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
	}
}
