package queue

//TypeNodeName 分类节点名
const TypeNodeName = "queue"

//Queue 消息队列配置
type Queue struct {
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}
