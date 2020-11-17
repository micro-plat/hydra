package proxy

import (
	"fmt"
	"net/url"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var clusters = cmap.New(2)

type UpCluster struct {
	c    conf.ICluster
	name string
}

//Next 获取下一个可用的上游地址
func (c *UpCluster) Next() (u *url.URL, err error) {

	node, ok := c.c.Next()
	if !ok {
		return nil, fmt.Errorf("集群%s无可用服务器", c.name)
	}
	path := fmt.Sprintf("http://%s:%s", node.GetHost(), node.GetPort())
	url, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("集群的服务器地址不合法:%s %s", path, err)
	}
	return url, nil
}
