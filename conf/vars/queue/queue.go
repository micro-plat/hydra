package queue

import "github.com/micro-plat/hydra/conf/pkgs/security"

//TypeNodeName 分类节点名
const TypeNodeName = "queue"

//Queue 消息队列配置
type Queue struct {
	security.ConfEncrypt
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}
