package queue

//VarRootName 在var中的跟路径
const VarRootName = "queue"

//Queue 消息队列配置
type Queue struct {
	Proto string `json:"proto,omitempty"`
	Raw   []byte `json:"-"`
}
