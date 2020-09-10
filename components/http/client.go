package http

import (
	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/hydra/components/pkgs/http"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

const (
	//typeNode DB在var配置中的类型名称
	dbTypeNode = "http"

	//nameNode DB名称在var配置中的末节点名称
	dbNameNode = "client"
)

//StandardHTTPClient db
type StandardHTTPClient struct {
	c container.IContainer
}

//NewStandardHTTPClient 创建DB
func NewStandardHTTPClient(c container.IContainer) *StandardHTTPClient {
	return &StandardHTTPClient{c: c}
}

//GetRegularClient 获取正式的没有异常数据库实例
func (s *StandardHTTPClient) GetRegularClient(names ...string) (d IClient) {
	d, err := s.GetClient(names...)
	if err != nil {
		panic(err)
	}
	return d
}

//GetClient 获取http请求对象
func (s *StandardHTTPClient) GetClient(names ...string) (d IClient, err error) {
	name := types.GetStringByIndex(names, 0, dbNameNode)
	obj, err := s.c.GetOrCreate(dbTypeNode, name, func(js *conf.RawConf) (interface{}, error) {
		ctx := context.Current()
		opt := []http.Option{
			http.WithRequestID(ctx.User().GetRequestID()),
		}
		if js == nil {
			return http.NewClient(opt...)
		}
		raw, err := http.WithRaw(js.GetRaw())
		if err != nil {
			return nil, err
		}
		opt = append(opt, raw)
		return http.NewClient(opt...)
	})
	if err != nil {
		return nil, err
	}
	return obj.(IClient), nil
}
